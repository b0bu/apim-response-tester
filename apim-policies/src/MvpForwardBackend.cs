using System.Text.RegularExpressions;
using Microsoft.Azure.ApiManagement.PolicyToolkit.Authoring;
using Microsoft.Azure.ApiManagement.PolicyToolkit.Authoring.Expressions;

namespace Mvp.Apis.Policies;

[Document]
public class MvpForwardBackend : IDocument
{
    private static string IdFromUrl(IExpressionContext context) =>
        Regex.Match(context.Request.OriginalUrl.ToString(), @"/job/(\d+)").Groups[1].Value;

    // get header, get value, get job id
    private static string IdFromHeader(IExpressionContext context) =>
        Regex.Match(context.Response.Headers.GetValueOrDefault("operation-location", ""), @"/job/(\d+)", RegexOptions.IgnoreCase).Groups[1].Value;

    private static bool IsGet(IExpressionContext context) =>
        context.Request.Method.Equals("GET");
    
    // ensure it's anything but a GET
    private static bool ShouldCache(IExpressionContext context) =>
        !context.Request.Method.Equals("GET");

    private static bool HasCachedResponse(IExpressionContext context) =>
        context.Variables.ContainsKey("cachedResponse");

    private static string CachedResponseBaseUrl(IExpressionContext context) =>
        (string)context.Variables["cachedResponse"];

    //  this sets oploc based on oploc which doesn't exist
    //  job and id need to be concat 
    //  nothing has been cached yet so can't look up the var

    private static string OperationLocationAbsolute(IExpressionContext context) {
        // on get op location header isn't present, so fall back to check cache
        var op = context.Response.Headers.GetValueOrDefault("operation-location", "");
        if (!string.IsNullOrEmpty(op))
        {
            // when POST operation-location is set by backend
            return string.Concat("https://policy-testing.azure-api.net/api/v1",
                    Regex.Match(op, @"/job/\d+").Value);
        } else
        {
            // when GET original url has the id in the url
            return string.Concat("https://policy-testing.azure-api.net/api/v1/job/",
                    (string)context.Variables["cachedResponse"]);
        }
    }

    private static string OpLocRequestUrl(IExpressionContext context) {
        var opLoc = context.Response.Headers.GetValueOrDefault("operation-location", "");
        return Regex.Replace(opLoc, @"/job/\d+.*$", string.Empty, RegexOptions.IgnoreCase);
    }

    public void Inbound(IInboundContext context)
    {
        context.Base();

        if (IsGet(context.ExpressionContext))
        {
            // only lookup on GET
            context.CacheLookupValue(new CacheLookupValueConfig {
                Key = IdFromUrl(context.ExpressionContext),
                VariableName = "cachedResponse",
                CachingType = "internal"
            });
            if (HasCachedResponse(context.ExpressionContext))
            {
                context.SetBackendService(new SetBackendServiceConfig {
                    BaseUrl = CachedResponseBaseUrl(context.ExpressionContext)
                });
            }
        } else {
            context.SetBackendService(new SetBackendServiceConfig { BackendId = "apim-rt-pool" });
        }
    }

    public void Backend(IBackendContext context) {
        context.ForwardRequest();
    }

    public void Outbound(IOutboundContext context)
    {
        context.Base();

        // it's not GET (right now)
        if (ShouldCache(context.ExpressionContext)) {

            //TODO: check if in cache first

            context.CacheStoreValue(new CacheStoreValueConfig {
                Key = IdFromHeader(context.ExpressionContext),
                Value = OpLocRequestUrl(context.ExpressionContext),
                Duration = 1200,
                CachingType = "internal"
            });
        }

        context.SetHeader("operation-location", OperationLocationAbsolute(context.ExpressionContext));

    }

    public void OnError(IOnErrorContext context)
    {
        context.Base();
        context.SetHeader("X-Error", "An error occurred in the API pipeline TRACE for more information");
    }
}

