// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
import PocketBase, { type AuthModel } from 'pocketbase';

declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface Platform {}
		// src/app.d.ts
		interface Locals {
			pb: PocketBase;
			user: AuthModel;
		}
	}
}

export {};
