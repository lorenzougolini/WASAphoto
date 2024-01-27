import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import SessionLoginView from '../views/SessionLoginView.vue'
import ProfileView from '../views/ProfileView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/stream', component: HomeView},
		{path: '/session', component: SessionLoginView},
		{path: '/users/:username', component: ProfileView}
	]
})

export default router
