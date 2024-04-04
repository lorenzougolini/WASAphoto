import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/stream', component: HomeView},
		{path: '/session', component: LoginView},
		{path: '/profile', component: ProfileView},
		{path: '/users/:username', component: ProfileView},
	]
});

export default router
