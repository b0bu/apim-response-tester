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

    private static bool HasCachedResponse(IExpressionContext context) =>
        context.Variables.ContainsKey("cachedResponse");

    private static string CachedResponseBaseUrl(IExpressionContext context) =>
        (string)context.Variables["cachedResponse"];

    private static string OperationLocationAbsolute(IExpressionContext context) =>
        string.Concat("https://policy-testing.azure-api.net/api/v1",
            Regex.Match(
                context.Response.Headers.GetValueOrDefault("operation-location", ""),
                @"/job/\d+").Value);

    private static string CurrentRequestUrl(IExpressionContext context) =>
        context.Response.Headers.GetValueOrDefault("operation-location", "");

    // ---------- Sections ----------
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

        // check if in cache first
        
        context.CacheStoreValue(new CacheStoreValueConfig {
            Key = IdFromHeader(context.ExpressionContext),
            Value = CurrentRequestUrl(context.ExpressionContext),
            Duration = 600,
            CachingType = "internal"
        });

        context.SetHeader("operation-location", OperationLocationAbsolute(context.ExpressionContext));
    }

    public void OnError(IOnErrorContext context)
    {
        context.Base();
        context.SetHeader("X-Error", "An error occurred in the API pipeline.");
    }
}

