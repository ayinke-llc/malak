import { CodeExamples } from './CodeExamples';

interface ApiExamplesProps {
  chartReference: string;
  integrationReference: string;
}

export function ApiExamples({ chartReference, integrationReference }: ApiExamplesProps) {
  const curlExample = `curl -H "Authorization: Bearer YOUR_TOKEN" \\
  -H "Content-Type: application/json" \\
  -X POST \\
  -d '{"value": 100}' \\
  https://api.malak.vc/v1/workspaces/integrations/${integrationReference}/charts/${chartReference}/points`;

  const typescriptExample = `import axios from 'axios';

async function addDataPoint(token: string, value: number) {
  const response = await axios.post(
    'https://api.malak.vc/v1/workspaces/integrations/${integrationReference}/charts/${chartReference}/points',
    { value },
    {
      headers: {
        'Authorization': \`Bearer \${token}\`,
        'Content-Type': 'application/json',
      },
    }
  );
  return response.data;
}

// Example usage:
// addDataPoint('your_token_here', 100)`;

  const pythonExample = `import requests

def add_data_point(token: str, value: int) -> dict:
    response = requests.post(
        'https://api.malak.vc/v1/workspaces/integrations/${integrationReference}/charts/${chartReference}/points',
        json={'value': value},
        headers={
            'Authorization': f'Bearer {token}',
            'Content-Type': 'application/json',
        }
    )
    return response.json()

# Example usage:
# add_data_point('your_token_here', 100)`;

  const goExample = `package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

func addDataPoint(token string, value int) error {
    payload := map[string]interface{}{
        "value": value,
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    req, err := http.NewRequest(
        "POST",
        "https://api.malak.vc/v1/workspaces/integrations/${integrationReference}/charts/${chartReference}/points",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return err
    }

    req.Header.Set("Authorization", "Bearer " + token)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return nil
}

// Example usage:
// addDataPoint("your_token_here", 100)`;

  const examples = [
    {
      language: 'shell',
      code: curlExample,
      label: 'cURL',
    },
    {
      language: 'typescript',
      code: typescriptExample,
      label: 'TypeScript',
    },
    {
      language: 'python',
      code: pythonExample,
      label: 'Python',
    },
    {
      language: 'go',
      code: goExample,
      label: 'Go',
    },
  ];

  return (
    <div className="w-full">
      <h2 className="text-2xl font-bold mb-4">API Examples</h2>
      <p className="text-muted-foreground mb-6">
        Here are examples of how to add data points to your chart using different programming languages.
        The examples show how to add a value of 100 to your chart.
      </p>
      <CodeExamples examples={examples} />
    </div>
  );
} 
