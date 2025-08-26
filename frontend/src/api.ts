const BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8080";


export async function register(email: string, password: string, displayName: string) {
    const res = await fetch(`${BASE}/api/register`, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ email, password, displayName }) });
    if (!res.ok) throw new Error((await res.json()).error || "register failed");
    return res.json();
}


export async function login(email: string, password: string) {
    const res = await fetch(`${BASE}/api/login`, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ email, password }) });
    if (!res.ok) throw new Error((await res.json()).error || "login failed");
    return res.json() as Promise<{ token: string; displayName: string }>
}


export function authHeader(): HeadersInit | undefined {
    const t = localStorage.getItem("token");
    return t ? { Authorization: `Bearer ${t}` } : undefined;
}


export async function me() {
    const res = await fetch(`${BASE}/api/me`, { headers: authHeader() });
    if (!res.ok) throw new Error("unauthorized");
    return res.json();
}


export async function getMessages(limit = 50) {
    const res = await fetch(`${BASE}/api/messages?limit=${limit}`, { headers: authHeader() });
    if (!res.ok) throw new Error("failed");
    return res.json();
}


export function wsUrl(token: string) {
    const url = new URL(BASE);
    url.protocol = url.protocol.replace("http", "ws");
    url.pathname = "/ws";
    url.searchParams.set("token", token);
    return url.toString();
}