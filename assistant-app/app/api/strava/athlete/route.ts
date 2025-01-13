import { getTokenFromCookies } from '../../../../utils/auth';

export async function GET(request) {
    const cookieHeader = request.headers.get('cookie');
    const accessToken = getTokenFromCookies(cookieHeader);

    if (!accessToken) {
        return new Response(JSON.stringify({ message: "No token found" }), {
            status: 401,
        });
    }

    try {
        const response = await fetch('https://www.strava.com/api/v3/athlete', {
            headers: {
                'Authorization': `Bearer ${accessToken}`,
            },
        });

        if (!response.ok) {
            throw new Error('Failed to fetch athlete info');
        }

        const data = await response.json();
        const headers = new Headers();
        headers.append(
            "Set-Cookie",
            `athleteId=${data.id}; HttpOnly; Secure; Path=/;`
        );

        return new Response(JSON.stringify(data), {
            status: 200,
            headers: headers,
        });
    } catch (error) {
        return new Response(JSON.stringify({ message: error.message }), {
            status: 500,
        });
    }
}