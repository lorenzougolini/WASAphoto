<script>
import PhotoCard from '../components/PhotoCard.vue';
import SearchBar from '../components/SearchBar.vue';

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
			
			photos: null,
			emptyStreamBanner: false,
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

				this.emptyStreamBanner = false;
				this.photos = response.data.Posts;
				if (this.photos) {
					this.photos.sort((a, b) => b.DateAndTime - a.DateAndTime);
				} else {
					this.emptyStreamBanner = true;
				}

			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},

		updatePhoto(updatedPhoto) {
            const idx = this.photos.findIndex((photo) => photo.PhotoID === updatedPhoto.PhotoID);
            if (idx !== -1) {
                this.photos[idx] = updatedPhoto;
            }
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
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="loadStream()">
						Refresh
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		<div v-show="emptyStreamBanner" class="alert alert-warning" role="alert">
			Wow, so empty! Start following users!
		</div>
		<div>
			<LoadingSpinner :loading="loading" />
		</div>
		<div class="stream-container">
			<div v-for="photo in this.photos" :key="photo.PhotoID">
				<PhotoCard 
					:photo="photo"
					@update-photo="updatePhoto"
					/>
				<br>
			</div>
		</div>
		
	</div>
</template>

<style>
</style>
