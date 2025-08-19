# populated from .env
#subscriptionId=""
#resourceGroup=""
serviceName="policy-testing"
apiName="apim-response-tester"
accessBearer=$(az account get-access-token | jq .accessToken | tr -d '"')

apiId=$(az apim api list \
  --resource-group "${resourceGroup}" \
  --service-name "${serviceName}" \
  --filter-display-name "${apiName}" \
  --query "[].id" \
  -o tsv)

az rest --method post \
 --uri "https://management.azure.com/subscriptions/${subscriptionId}/resourceGroups/${resourceGroup}/providers/Microsoft.ApiManagement/service/${serviceName}/gateways/managed/listDebugCredentials?api-version=2023-05-01-preview" \
 --body "{
    \"credentialsExpireAfter\": \"PT1H\",
    \"apiId\": \"${apiId}\",
    \"purposes\": [\"tracing\"]
}" \
    --query token -o tsv
