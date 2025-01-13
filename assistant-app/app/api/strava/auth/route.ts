import axios from 'axios';

// Get the Auth token from Strava

const stravaClientId = process.env.STRAVA_CLIENT_ID;
const stravaClientSecret = process.env.STRAVA_CLIENT_SECRET;

export async function POST(request) {
    const { code } = await request.json();
    console.log("Auth code received", code);

    try {
        const response = await axios.post<{ access_token: string; expires_in: number }>(
            "https://www.strava.com/oauth/token",
            {
                client_id: stravaClientId,
                client_secret: stravaClientSecret,
                code: code,
                grant_type: "authorization_code",
            }
        );

        // Set Access Token as HTTP-Only Cookie
        const headers = new Headers();
        headers.append(
            "Set-Cookie",
            `stravaAccessToken=${response.data.access_token}; HttpOnly; Secure; Path=/; Max-Age=${response.data.expires_in}`
        );
        console.log("Token generated, ", response.data.access_token);
        return new Response(JSON.stringify({ message: "Token generated", accessToken: response.data.access_token }), {
            status: 200,
            headers: headers,
        });
    } catch (error) {
        console.error("Error exchanging code for token:", error);
        return new Response(JSON.stringify({ message: "Failed to exchange code for token" }), {
            status: 500,
        });
    }
}

