import { getTokenFromCookies } from '../../../../utils/auth';

export async function GET(request) {
    const cookieHeader = request.headers.get('cookie');
    const accessToken = getTokenFromCookies(cookieHeader);

    if (accessToken) {
        return new Response(JSON.stringify({ accessToken: accessToken }), {
            status: 200,
        });
    } else {
        return new Response(JSON.stringify({ message: "No token found" }), {
            status: 404,
        });
    }
}
