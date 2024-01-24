<script>
import Photo from '../components/Photo.vue';
import {user} from '../stores/user.js';

export default {
    components: {
        Photo,
    },
    
	data: function() {
		return {
			errormsg: null,
			loading: false,
            username: sessionStorage.getItem("username"),
            userid: sessionStorage.getItem("userid"),
            shownUsername: this.$route.params.username,
            profileJson: {},
            picture: null,
		}
	},

	methods: {
		async loadProfile (shownUsername) {
			this.loading = true;
			this.errormsg = null;

            if (sessionStorage.getItem("authenticated")) {
                
                try {

                    console.log(this.username, this.userid);

                    let response = await this.$axios.get("/users/" + shownUsername, {
                        headers: {
                            'Authorization': this.userid,
                        }
                    });
                    this.profileJson = response.data;

                } catch (e) {
                    this.errormsg = e.toString();
                }
                this.loading = false;

            } else {
                this.errormsg = "You are not logged in!";
            }
		},

        fileUpload(event) {
            console.log(event.target.files[0]);
            this.picture = event.target.files[0];
        },

        async newPost(description, picture) {

            // if (sessionStorage.getItem("authenticated") && sessionStorage.getItem("username") == shownUsername) {
            // maybe the backend will do it for me without checking here

            const bodyFormData = new FormData();
            bodyFormData.append('description', description);
            bodyFormData.append('picture', picture);
            
            this.loading = true;   
            this.errormsg = null;

            try {         
                await this.$axios.post("/users/" + this.username + "/photos/", bodyFormData, {
                    headers: {
                        'Authorization': this.userid,
                        "Content-Type": "multipart/form-data"
                    }
                });
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            console.log("New post created!");
            this.loadProfile(this.username);
        },

        async deletePost(photoID) {
            this.loading = true;   
            this.errormsg = null;

            try {         
                await this.$axios.delete("/users/" + this.username + "/photos/" + photoID, {
                    headers: {
                        'Authorization': this.userid,
                    }
                });
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            console.log("Post deleted!");
            this.loadProfile(this.username);
        },
	},

    computed: {
        numberOfPosts() {
            if (this.profileJson.Posts){
                return this.profileJson.Posts.NumberOfPosts
            } else {
                return 0;
            }
        },
        numberOfFollowers() {
            if (this.profileJson.Posts){
                return this.profileJson.Followers.NumberOfFollowers
            } else {
                return 0;
            }
        },
        numberOfFollowing() {
            if (this.profileJson.Posts){
                return this.profileJson.Followers.NumberOfFollowing
            } else {
                return 0;
            }
        },
        profilePhotos() {
            if (this.profileJson.Posts){
                return this.profileJson.Posts.Photos
            } else {
                return [];
            }
        }
    },

	mounted() {
		this.loadProfile(this.shownUsername)
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Profile page</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="loadProfile(shownUsername)">
						Refresh
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        <div class="profile-container">
            <div class="photo-container">
                <h2 class="h2">This is the profile of {{ shownUsername }}</h2>
                <p>Number of posts: {{ numberOfPosts }}, 
                    Number of followers: {{ numberOfFollowers }}, 
                    Number of following: {{ numberOfFollowing }}
                </p>
                <div v-if="this.profileJson.Posts" class="horizontal-photo-container">
                    <div v-for="photo in this.profilePhotos" :key="photo.PhotoID" class="horizontal-photo-div">
                        <!-- {{ photo.PhotoID }}, {{ photo.Description }}, {{ photo.DateAndTime }}
                        <img v-bind:src="`/pictures/${photo.PhotoID}.jpg`" alt="{{photo.Description}}"> -->
                        <Photo @delete-post="deletePost(photo.PhotoID)"
                            :photoLocation="`/pictures/${photo.PhotoID}.jpg`"
                            :photoDescription="photo.Description"
                            :photoDate="photo.DateAndTime"
                        />
                            <!-- :photoId="photo.PhotoID" -->

                        </div>
                </div>
            </div>

            <div class="new-post-container">
                <h3>New Post</h3>
                <form @submit.prevent="newPost">
                    Description: <input type="text" v-model="description" /><br />
                    Picture: <input type="file" v-on:change="fileUpload" /><br />
                    <br>
                    <div class="btn-group me-2">
                        <button type="button" class="btn btn-sm btn-outline-primary" @click="newPost(description, picture)">
                            New Post
                        </button>
                    </div>
                </form>
            </div>
        </div>
	</div>
</template>

<style>
/* maybe do it with a grid */
.horizontal-photo-container {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
}
.horizontal-photo-div {
    margin-right: 10px;
}

.profile-container {
    display: flex;
    flex-direction: row;
}
.photo-container {
    width: 80%;
}
.new-post-container {
    width: 30%;
}
</style>
