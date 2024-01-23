<script>
import ProfileView from './ProfileView.vue'
import {user} from '../stores/user.js';

export default {


	data: function() {
		return {
			errormsg: null,
			loading: false,
			username: "",
			logged: {},
		}
	},

	methods: {
		async login() {
			this.loading = true;
			this.errormsg = null;

			try {
				let response = await this.$axios.post("/session?username=" + this.username);
				
				user.value.userid = response.data.UserID;
				user.value.username = response.data.Username;
				user.value.authenticated = true;
	
				sessionStorage.setItem("userid", user.value.userid);
				sessionStorage.setItem("username", user.value.username);
				sessionStorage.setItem("authenticated", user.value.authenticated);
				// console.log(sessionStorage.userid, sessionStorage.username, sessionStorage.authenticated);
				
				this.$router.replace("/users/" + user.value.username, { 
					headers: {
						'Authorization': user.value.userid,
					}
				});
				
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Login</h1>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<form @submit.prevent="login">
			Username: <input type="text" v-model="username" /><br />
			<button class="btn-group me-2" type="submit">
				Login
			</button>
		</form>
	</div>
</template>

<style>
</style>
