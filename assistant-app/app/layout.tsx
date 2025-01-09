import { Inter } from "next/font/google";
import "./globals.css";
import Warnings from "./components/warnings";
import { assistantId } from "./assistant-config";
import { Providers } from "./providers";
import Navigation from "./navbar";
const inter = Inter({ subsets: ["latin"] });

export const metadata = {
  title: "Cycling Assistant",
  description: "A quickstart template using the Assistants API with OpenAI",
  icons: {
    icon: "/openai.svg",
  },
};

import { ReactNode } from "react";

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en" >
      <body className={inter.className}>
      <Navigation />
        <Providers>
          {assistantId ? children : <Warnings />}
          <img className="logo" src="/openai.svg" alt="OpenAI Logo" />
        </Providers>
      </body>
    </html>
  );
}