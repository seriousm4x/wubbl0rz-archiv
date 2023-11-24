import { pb } from '$lib/pocketbase';
import { dev } from '$app/environment';

pb.authStore.loadFromCookie(document.cookie);
pb.authStore.onChange(() => {
	document.cookie = pb.authStore.exportToCookie({ httpOnly: false, secure: dev ? false : true });
}, true);
