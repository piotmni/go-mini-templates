import { betterAuth } from "better-auth";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { openAPI, jwt, deviceAuthorization } from "better-auth/plugins";
import { db } from "./db";
import { v7 as uuidv7 } from "uuid";

export const auth = betterAuth({
  database: drizzleAdapter(db, {
    provider: "pg",
  }),

  advanced: {
    database: {
      generateId: () => uuidv7(),
    },
  },

  socialProviders: {
    github: {
      clientId: process.env.GITHUB_CLIENT_ID || "",
      clientSecret: process.env.GITHUB_CLIENT_SECRET || "",
    },
  },

  trustedOrigins: ["http://localhost:3000"],

  plugins: [
    openAPI(),
    deviceAuthorization({
      verificationUri: "/device",
    }),
    jwt({
      jwt: {
        definePayload: ({ user }) => ({
          id: user.id,
          email: user.email,
        }),
      },
    }),
  ],
});
