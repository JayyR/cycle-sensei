import { getTokenFromCookies } from "@/utils/auth";

export async function fetchWithAuth(url: string, request: Request) {
    const cookieHeader = request.headers.get('cookie');
    const stravaToken = getTokenFromCookies(cookieHeader);
    const response = await fetch(url, {
        headers: {
            'StravaAuthToken': `${stravaToken}`
        }
    });
    if (!response.ok) {
        throw new Error('Network response was not ok');
    }
    return response.json();
}
