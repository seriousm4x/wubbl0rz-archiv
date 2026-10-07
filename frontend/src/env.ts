import { defineEnvVars } from '@sveltejs/kit/env';

export const variables = defineEnvVars({
	PRIVATE_ALLOW_SEARCH_INDEXING: {},
	PRIVATE_MEILI_ADMIN_KEY: {},
	PUBLIC_API_URL: { public: true },
	PUBLIC_FRONTEND_URL: { public: true },
	PUBLIC_MEILI_SEARCH_KEY: { public: true },
	PUBLIC_MEILI_URL: { public: true }
});
