using System.Text.RegularExpressions;
using Microsoft.Azure.ApiManagement.PolicyToolkit.Authoring;
using Microsoft.Azure.ApiManagement.PolicyToolkit.Authoring.Expressions;

namespace Mvp.Apis.Policies;

[Document]
public class MvpForwardBackend : IDocument
{
    // used on /GET to parse out id finds 'job' gets id
    [Expression]
    private static string JobID(IExpressionContext ctx)
        => Regex.Match(ctx.Request.OriginalUrl.ToString(), @"/job/\d+").Value;

    public void Inbound(IInboundContext context)
    {
        context.Base();
        context.SetBackendService(new SetBackendServiceConfig { BackendId = "apim-response-tester" });
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

    public void Outbound(IOutboundContext context)
    {
        context.Base();
        context.SetHeader("operation-location", NormalizeHeader(context.ExpressionContext));
    }

    public void OnError(IOnErrorContext context)
    {
        context.Base();
        context.SetHeader("X-Error", "An error occurred in the API pipeline.");
    }
}

