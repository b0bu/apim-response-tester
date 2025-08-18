using Microsoft.Azure.ApiManagement.PolicyToolkit.Authoring;
using Microsoft.Azure.ApiManagement.PolicyToolkit.Authoring.Expressions;

namespace Mvp.Apis.Policies;

[Document]
public class MvpForwardBackend : IDocument
{
    public void Inbound(IInboundContext context)
    {
        context.SetHeader("operation-location", "@(context.Response.Headers.GetValueOrDefault(\"operation-location\", \"\").Replace(\"apim-response-tester-0.uksouth.azurecontainer.io\", \"https://policy-testing.azure-api.net/api/v1\"))");
    }

    public void Backend(IBackendContext context)
    {
        context.SetBackendService(new SetBackendServiceConfig { BaseUrl = "http://apim-response-tester.dmf6h7gqcrcef7b9.uksouth.azurecontainer.io:8080" });
    }

    public void Outbound(IOutboundContext context)
    {
    }

    public void OnError(IOnErrorContext context)
    {
        context.Base();
        context.SetHeader("X-Error", "An error occurred in the API pipeline.");
    }
}

