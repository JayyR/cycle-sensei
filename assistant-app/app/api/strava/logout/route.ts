export async function POST() {
    const headers = new Headers();
    headers.append(
        "Set-Cookie",
        `stravaAccessToken=; HttpOnly; Secure; Path=/; Max-Age=0`
    );

    // Clear athlete cache
    headers.append(
        "Clear-Site-Data",
        `"storage"`
    );

    // Clear athlete ID from session
    headers.append(
        "Set-Cookie",
        `athleteId=; HttpOnly; Secure; Path=/; Max-Age=0`
    );

    return new Response(JSON.stringify({ message: "Logged out" }), {
        status: 200,
        headers: headers,
    });
}
