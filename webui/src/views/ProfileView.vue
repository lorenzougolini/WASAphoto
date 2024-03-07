<script>
import PhotoCard from '../components/PhotoCard.vue';
import LoadingSpinner from '../components/LoadingSpinner.vue';
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
            shownUsername: this.$route.params.username,
            profileJson: {},
            picture: null,
		}
	},

	methods: {
		async loadProfile(shownUsername) {
			this.loading = true;
			this.errormsg = null;

            if (sessionStorage.getItem("authenticated")) {
                
                try {

                    // console.log(this.username, this.userid);
                    let response = await this.$axios.get("/users/" + shownUsername, {
                        headers: {
                            'Authorization': this.userid,
                        }
                    });
                    this.profileJson = response.data;

                } catch (e) {
                    this.errormsg = e.toString();
                    this.$router.go(-1);
                }
                this.loading = false;

            } else {
                this.errormsg = "You are not logged in!";
            }
		},

        async newPost(description, picture) {
            
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

        async followUser(){
            this.loading = true;   
            this.errormsg = null;
            
            try {        
                await this.$axios.put("/users/" + this.username + "/followers/" + this.shownUsername, null, {
                    headers: {
                        'Authorization': this.userid,
                    }
                });
                console.log("User followed!");
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            this.loadProfile(this.shownUsername);
        },

        async unfollowUser() {
            this.loading = true;   
            this.errormsg = null;
            
            try {        
                await this.$axios.delete("/users/" + this.username + "/followers/" + this.shownUsername, {
                    headers: {
                        'Authorization': this.userid,
                    }
                });
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            console.log("User unfollowed!");
            this.loadProfile(this.shownUsername);            
        },
        
        fileUpload(event) {
            console.log(event.target.files[0]);
            this.picture = event.target.files[0];
        },

        followedByYou() {
            try {
                var followed = this.profileJson.Followers.Usernames.includes(this.username);
            } catch (e) {
                var followed = false;
            }
            return followed;
        },
	},
    
    computed: {
        numberOfPosts() {
            if (this.profileJson.Posts){
                return this.profileJson.Posts.length
            } else {
                return 0;
            }
        },
        numberOfFollowers() {
            if (this.profileJson.Followers){
                return this.profileJson.Followers.NumberOfFollowers
            } else {
                return 0;
            }
        },
        numberOfFollowing() {
            if (this.profileJson.Following){
                return this.profileJson.Following.NumberOfFollowing
            } else {
                return 0;
            }
        },
        profilePhotos() {
            if (this.profileJson.Posts){
                return this.profileJson.Posts
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
            <div class="d-flex align-items-center">
				<SearchBar @search-user="loadProfile(this.shownUsername)"/>
			</div>
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
            <LoadingSpinner :loading="loading" />
            <div v-if="!loading" class="photo-container">
                <h2 class="h2">This is the profile of {{ shownUsername }}</h2>
                <p>Posts: {{ numberOfPosts }}, 
                    Followers: {{ numberOfFollowers }}, 
                    Following: {{ numberOfFollowing }}
                </p>
                <div v-if="this.profileJson.Posts" class="horizontal-photo-container">
                    <div v-for="photo in this.profilePhotos" :key="photo.PhotoID" class="horizontal-photo-div">

                        <PhotoCard @delete-post="deletePost(photo.PhotoID)"
                            :photoAuthor="photo.Author"
                            :photoLocation="`/pictures/${photo.PhotoID}.jpg`"
                            :photoDescription="photo.Description"
                            :photoDate="photo.DateAndTime"
                            :photoLikes="photo.Likes"
                            :photoComments="photo.Comments"
                            :parent="profile"
                        />

                        </div>
                </div>
            </div>

            <div v-if="shownUsername == this.username" class="new-post-container">
                <h3>New Post</h3>
                <form @submit.prevent="newPost">
                    Picture: <input type="file" v-on:change="fileUpload" /><br />
                    Description: <input type="text" v-model="description" /><br />
                    <br>
                    <div class="btn-group me-2">
                        <button type="button" class="btn btn-sm btn-outline-primary" @click="newPost(description, picture)">
                            New Post
                        </button>
                    </div>
                </form>
            </div>
            <div v-else>
                <div v-if="this.profileJson.Followers">
                    <div v-if="followedByYou()">
                        <button class="btn btn-sm btn-outline-primary" @click="unfollowUser()">Unfollow</button>
                    </div>
                    <div v-else>
                        <button class="btn btn-sm btn-outline-primary" @click="followUser()">Follow</button>
                    </div>
                </div>
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
