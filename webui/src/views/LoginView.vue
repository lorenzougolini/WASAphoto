<script>
export default {

	data: function() {
		return {
			errormsg: null,
			loading: false,
			loggingUsername: "",
			logged: {},
		}
	},

	methods: {
		async login() {
			this.loading = true;
			this.errormsg = null; 
			
			try {
				let response = await this.$axios.post(`/session?username=${this.loggingUsername}`);
				
				sessionStorage.setItem("userid", response.data.UserID),
				sessionStorage.setItem("username", response.data.Username),
				sessionStorage.setItem("logged", true),
				
				this.$router.replace("/stream", { 
					headers: {
						'Authorization': sessionStorage.getItem("userid"),
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
			Username: <input type="text" v-model="loggingUsername" /><br />
			<button class="btn-group me-2" type="submit">
				Login
			</button>
		</form>
	</div>
</template>

<style>
</style>
