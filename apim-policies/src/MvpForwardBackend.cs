using System.Text.RegularExpressions;
using Microsoft.Azure.ApiManagement.PolicyToolkit.Authoring;
using Microsoft.Azure.ApiManagement.PolicyToolkit.Authoring.Expressions;

namespace Mvp.Apis.Policies;

[Document]
public class MvpForwardBackend : IDocument
{
    // used on /GET to parse out id finds 'job' gets id
    private static string JobID(string url)
        => Regex.Match(url, @"/job/\d+").Value;

    private static string ExtractJobId(IExpressionContext context) 
        => JobID(context.Request.OriginalUrl.ToString());

    private static string ExtractRequestMethod(IExpressionContext context)
        => context.Request.Method;

    private static string CachedResponse(IExpressionContext context, string key)
        => context.Variables[key].ToString();

    public void Inbound(IInboundContext context)
    {
        context.Base();
        // inbound context has no property Request
        context.CacheLookupValue(new CacheLookupValueConfig { Key = ExtractJobId(context.ExpressionContext),  VariableName = "cachedResponse", CachingType = "internal" });

        // if GET and cached use it, else balance
        if (ExtractRequestMethod(context.ExpressionContext) == "GET" && CachedResponse(context.ExpressionContext, "cachedResponse") is not null) {
            context.SetBackendService(new SetBackendServiceConfig { BaseUrl = CachedResponse(context.ExpressionContext, "cachedResponse") });
        } else {
            context.SetBackendService(new SetBackendServiceConfig { BackendId = "apim-rt-pool" });
        }
    }

    public void Backend(IBackendContext context)
    {
        context.ForwardRequest();
    }

    // match group [a-zA-Z0-9] if alphanumeric
    [Expression]
    private static string NormalizeHeader(IExpressionContext context) {
        var match = Regex.Match(context.Response.Headers.GetValueOrDefault("operation-location", ""), @"/job/\d+");
        return "https://policy-testing.azure-api.net/api/v1" + match.Value;
    }

    private static void CacheJob(IOutboundContext context, IExpressionContext ctx) {
        var backend_id = ctx.Request.OriginalUrl.ToString();
        var job_id = JobID(ctx.Request.OriginalUrl.ToString());
        // not production caching type
        context.CacheStoreValue(new CacheStoreValueConfig { Key = job_id, Value = backend_id, Duration = 600, CachingType = "internal" }); 
    }

    public void Outbound(IOutboundContext context)
    {
        context.Base();
        CacheJob(context, context.ExpressionContext);
        context.SetHeader("operation-location", NormalizeHeader(context.ExpressionContext));
    }

    public void OnError(IOnErrorContext context)
    {
        context.Base();
        context.SetHeader("X-Error", "An error occurred in the API pipeline.");
    }
}

