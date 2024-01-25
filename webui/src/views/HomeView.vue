<script>
import SessionLoginView from './SessionLoginView.vue';
import Photo from '../components/Photo.vue';
import {user} from '../stores/user.js';

export default {
	components: {
		SessionLoginView,
		Photo,
	},

	data: function() {
		return {
			errormsg: null,
			loading: false,
			username: sessionStorage.getItem("username"),
            userid: sessionStorage.getItem("userid"),
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
		this.loadStream()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home page</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
				</div>
				<!-- <div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
						New
					</button>
				</div> -->
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div v-if="!this.username || !this.userid">
			<SessionLoginView />
		</div>
		<div>
			<LoadingSpinner :loading="loading" />
		</div>
		<div class="stream-container">
			<div v-for="photo in this.streamJson.Posts" :key="photo.PhotoID">
				<Photo 
					:photoAuthor="photo.Author"
					:photoLocation="`/pictures/${photo.PhotoID}.jpg`" 
					:photoDescription="photo.Description"
					:photoDate="photo.DateAndTime"
					:photoLikes="photo.Likes"
					:photoComments="photo.Comments"
					/>
			</div>
		</div>
		
	</div>
</template>

<style>
</style>
