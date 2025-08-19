using System.Text.RegularExpressions;
using Microsoft.Azure.ApiManagement.PolicyToolkit.Authoring;
using Microsoft.Azure.ApiManagement.PolicyToolkit.Authoring.Expressions;

namespace Mvp.Apis.Policies;

[Document]
public class MvpForwardBackend : IDocument
{
    private static string JobId(IExpressionContext context) =>
        Regex.Match(context.Request.OriginalUrl.ToString(), @"/job/\d+").Value; // Regex.Match is allowed in APIM expressions

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
        context.Request.OriginalUrl.ToString();

    // ---------- Sections ----------
    public void Inbound(IInboundContext context)
    {
        context.Base();


        context.SetBackendService(new SetBackendServiceConfig { BackendId = "apim-rt-pool" });

        if (IsGet(context.ExpressionContext))
        {
            // only lookup on GET
            context.CacheLookupValue(new CacheLookupValueConfig {
                Key = JobId(context.ExpressionContext),
                VariableName = "cachedResponse",
                CachingType = "internal"
            });
            if (HasCachedResponse(context.ExpressionContext))
            {
                context.SetBackendService(new SetBackendServiceConfig {
                    BaseUrl = CachedResponseBaseUrl(context.ExpressionContext)
                });
            }
        }
    }

    public void Backend(IBackendContext context) {
        context.ForwardRequest();
    }

    public void Outbound(IOutboundContext context)
    {
        context.Base();

        context.CacheStoreValue(new CacheStoreValueConfig {
            Key = JobId(context.ExpressionContext),
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

