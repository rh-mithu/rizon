const BASE_URL = 'http://192.168.10.249:8080'; // Placeholder, should be ENV variable later

export const api = {
    requestLoginLink: async (email: string): Promise<void> => {
        try {
            const response = await fetch(`${BASE_URL}/api/v1/auth/request-link`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email: email }),
            });
            if (!response.ok) {
                let errorMessage = 'Failed to request login link';
                const text = await response.text();
                try {
                    const errorData = JSON.parse(text);
                    errorMessage = errorData.message || errorData.error || errorMessage;
                } catch (parseError) {
                    errorMessage = text || errorMessage;
                }
                throw new Error(errorMessage);
            }
        } catch (error: any) {
            console.error('API Error:', error);
            throw new Error(error.message || 'An unexpected error occurred');
        }
    },
};
