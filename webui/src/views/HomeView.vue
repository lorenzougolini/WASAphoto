<script>
import PhotoCard from '../components/PhotoCard.vue';
import SearchBar from '../components/SearchBar.vue';
import {user} from '../stores/user.js';

export default {
	components: {
		PhotoCard,
		SearchBar,
	},

	data: function() {
		return {
			errormsg: null,
			loading: false,
			username: sessionStorage.getItem("username"),
            userid: sessionStorage.getItem("userid"),
			logged: sessionStorage.getItem("logged"),
			streamJson: {},
		}
	},
	methods: {
		async loadStream() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/stream", {
					headers: {
						'Authorization': this.userid,
					}
				});

				this.streamJson = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
	mounted() {
		if (this.logged) {
			this.loadStream();
		} else {
			this.$router.replace("/session");
		}
	}
}
</script>

<template>
	<div >
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home page</h1>
			<div class="d-flex align-items-center">
				<SearchBar />
			</div>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div>
			<LoadingSpinner :loading="loading" />
		</div>
		<div class="stream-container">
			<div v-for="photo in this.streamJson.Posts" :key="photo.PhotoID">
				<PhotoCard 
					:photoAuthor="photo.Author"
					:photoId="photo.PhotoID"
					:photoDescription="photo.Description"
					:photoDate="photo.DateAndTime"
					:photoLikes="photo.Likes"
					:photoComments="photo.Comments"
					/>
				<br>
			</div>
		</div>
		
	</div>
</template>

<style>
</style>
