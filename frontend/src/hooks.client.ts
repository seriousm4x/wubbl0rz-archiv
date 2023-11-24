import { pb } from '$lib/pocketbase';
import { currentUser } from '$lib/stores/user';
import { dev } from '$app/environment';

pb.authStore.loadFromCookie(document.cookie);
pb.authStore.onChange(() => {
	currentUser.set(pb.authStore.model);
	document.cookie = pb.authStore.exportToCookie({ httpOnly: false, secure: dev ? false : true });
}, true);
