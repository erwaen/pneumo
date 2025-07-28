import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

interface APIResult {
  short_url: string;
}

const URLPneumifyForm: React.FC = () => {
  const [url, setUrl] = useState<string>('');
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [result, setResult] = useState<APIResult | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handlePneumifyClick = async (): Promise<void> => {
    setIsLoading(true);
    setResult(null);
    setError(null);

    try {
      const apiEndpoint = 'http://localhost:8081/pneumify/';


      const response = await fetch(apiEndpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url: url }),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: 'An unknown error occurred' }));
        throw new Error(errorData.message || 'The server responded with an error.');
      }

      const data: APIResult = await response.json();
      console.log('API response:', data);
      setResult(data);

    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full max-w-md flex flex-col items-center space-y-4">
      <div className="flex space-x-2 w-full">
        <Input
          type="url"
          placeholder="https://www.example.com"
          className="w-80 flex-grow"
          value={url}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => setUrl(e.target.value)}
          disabled={isLoading}
        />
        <Button onClick={handlePneumifyClick} disabled={isLoading}>
          {isLoading ? 'Pneumifying...' : 'Pneumify URL!'}
        </Button>
      </div>

      {/* success message */}
      {result && result.short_url && (
        <div className="p-3 bg-green-100 text-green-800 rounded-md w-full text-center dark:bg-green-900/30 dark:text-green-300">
          <p>
            Success! Your link is:{' '}
            <a href={result.short_url} target="_blank" rel="noopener noreferrer" className="font-bold underline hover:text-green-900 dark:hover:text-green-200">
              {result.short_url}
            </a>
          </p>
        </div>
      )}

      {/* error message */}
      {error && (
        <div className="p-3 bg-red-100 text-red-800 rounded-md w-full text-center dark:bg-red-900/30 dark:text-red-300">
          <p>{error}</p>
        </div>
      )}
    </div>
  );
};

export default URLPneumifyForm;

