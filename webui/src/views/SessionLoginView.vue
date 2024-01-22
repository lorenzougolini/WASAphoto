<script>
import ProfileView from './ProfileView.vue'
export default {
	
	components: {
		ProfileView
	},

	data: function() {
		return {
			errormsg: null,
			loading: false,
			logged: {},
		}
	},
	methods: {
		async login() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.post("/session?username=" + this.username);
				this.$router.push("/users/" + this.username, { // change to replace
					headers: {
						'Authorization': response.data.UserID
					}
				});
				this.logged = response.data;
				// this.logged[0] = response.data.UserID;
				// this.logged[1] = response.data.Username;
				console.log(this.logged);
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
	<!-- <div>
		<ProfileView :userid = "userid"/>
	</div> -->
</template>

<style>
</style>
