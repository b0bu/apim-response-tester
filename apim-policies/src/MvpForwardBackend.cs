using Microsoft.Azure.ApiManagement.PolicyToolkit.Authoring;
using Microsoft.Azure.ApiManagement.PolicyToolkit.Authoring.Expressions;

namespace Mvp.Apis.Policies;

[Document]
public class MvpForwardBackend : IDocument
{
    public void Inbound(IInboundContext context)
    {
        context.Base();
        context.SetBackendService(new SetBackendServiceConfig { BackendId = "apim-response-tester" });
    }

    public void Backend(IBackendContext context)
    {
        context.ForwardRequest();
    }

    public void Outbound(IOutboundContext context)
    {
        context.Base();
        context.SetHeader("operation-location", "https://policy-testing.azure-api.net/api/v1");
    }

    public void OnError(IOnErrorContext context)
    {
        context.Base();
        context.SetHeader("X-Error", "An error occurred in the API pipeline.");
    }
}

