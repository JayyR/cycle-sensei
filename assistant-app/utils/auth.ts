
export const getTokenFromCookies = (cookieHeader: string | null) => {
    if (!cookieHeader) {
        return null;
    }

    const token = cookieHeader.split('; ').find(row => row.startsWith('stravaAccessToken='));
    if (token) {
        return token.split('=')[1];
    }

    return null;
};