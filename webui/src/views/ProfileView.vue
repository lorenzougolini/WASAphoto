<script>
import PhotoCard from '../components/PhotoCard.vue';
import LoadingSpinner from '../components/LoadingSpinner.vue';
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
            
            isOwner: false,

            pathUser: this.$route.params.username,
            
            profileData: {},
            profilePhotos: [],
            userNotFound: false,
            
            description: '',
            newUsername: '',
            usernameTakenBanner: false,
		}
	},

	methods: {
		async loadProfile(pathUser) {
			this.loading = true;
			this.errormsg = null;
            try {

                let response = await this.$axios.get("/users/" + pathUser, {
                    headers: {
                        'Authorization': this.userid,
                    }
                });

                const profile = response.data;

                this.profileData = {
                    Username: profile.Username, 
                    Followers: profile.Followers, 
                    Following: profile.Following,
                    Banned: profile.Banned
                };

                this.profilePhotos = profile.Posts;
                
                if (this.profileData.Username === this.username) {
                    this.isOwner = true;
                } else {
                    this.isOwner = false;
                }

            } catch (e) {
                if (e.response && e.response.status === 404) {
                    this.userNotFound = true;                    
                } else {
                    this.errormsg = e.toString();
                }
            }
            this.loading = false;
		},

        async newPost(description, picture) {
            
            const bodyFormData = new FormData();
            bodyFormData.append('description', description);
            bodyFormData.append('picture', picture);
            
            this.loading = true;   
            this.errormsg = null;

            try {         
                await this.$axios.post("/users/" + this.username + "/photos", bodyFormData, {
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
            
            // Clear form fields
            document.getElementById('post-form').reset();
            this.description = '';
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
                
                console.log("Post deleted!");

            } catch (e) {
                this.errormsg = e.toString();
            }

            this.loading = false;
            this.loadProfile(this.username);
        },

        async changeUsername(newUsername) {
            this.loading = true;   
            this.errormsg = null;
            
            try {        
                let response = await this.$axios.put(`/users/${this.username}`, newUsername, {
                    headers: {
                        'Authorization': this.userid,
                    }
                });

                sessionStorage.setItem("username", response.data.Username);
                
                console.log("Username changed!");
                this.pathUser = response.data.Username;
                this.username = response.data.Username;
                this.$router.push("/users/" + response.data.Username);
                this.loadProfile(this.pathUser);

            } catch (e) {
                if (e.response && e.response.status === 400) {
                    this.usernameTakenBanner = true;
                    setTimeout(() => {
                        this.usernameTakenBanner = false;
                    }, 5000);
                    
                } else {
                    this.errormsg = e.toString();
                }
            }
            this.loading = false;
        },

        async followUser(){
            this.loading = true;   
            this.errormsg = null;
            
            try {        
                await this.$axios.put("/users/" + this.username + "/followers/" + this.pathUser, null, {
                    headers: {
                        'Authorization': this.userid,
                    }
                });
                
                console.log("User followed!");
            
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            this.loadProfile(this.pathUser);
        },

        async unfollowUser() {
            this.loading = true;   
            this.errormsg = null;
            
            try {        
                await this.$axios.delete("/users/" + this.username + "/followers/" + this.pathUser, {
                    headers: {
                        'Authorization': this.userid,
                    }
                });

                console.log("User unfollowed!");
            
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            this.loadProfile(this.pathUser);            
        },

        async banUser() {
            this.loading = true;   
            this.errormsg = null;
            
            try {        
                await this.$axios.put("/users/" + this.username + "/banned/" + this.pathUser, null, {
                    headers: {
                        'Authorization': this.userid,
                    }
                });
                
                console.log("User banned!");

            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            this.loadProfile(this.pathUser);
        },
        
        async unbanUser() {
            this.loading = true;   
            this.errormsg = null;
            
            try {        
                await this.$axios.delete("/users/" + this.username + "/banned/" + this.pathUser, {
                    headers: {
                        'Authorization': this.userid,
                    }
                });
                
                console.log("User unbanned!");
            
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            this.loadProfile(this.pathUser);
        },
        
        fileUpload(event) {
            console.log(event.target.files[0]);
            this.picture = event.target.files[0];
        },

        followedByYou() {
            try {
                var followed = this.profileData.Followers.Usernames.includes(this.username);
            } catch (e) {
                var followed = false;
            }
            return followed;
        },

        updatePhoto(updatedPhoto) {
            const idx = this.profilePhotos.findIndex((photo) => photo.PhotoID === updatedPhoto.PhotoID);
            if (idx !== -1) {
                this.profilePhotos[idx] = updatedPhoto;
            }
        }

	},
    
    computed: {
        numberOfPosts() {
            if (this.profilePhotos){
                return this.profilePhotos.length
            } else {
                return 0;
            }
        },
        numberOfFollowers() {
            if (this.profileData.Followers){
                return this.profileData.Followers.NumberOfFollowers
            } else {
                return 0;
            }
        },
        numberOfFollowing() {
            if (this.profileData.Following){
                return this.profileData.Following.NumberOfFollowing
            } else {
                return 0;
            }
        },
    },

    watch: {
        '$route'(to, from) {

            if (from.path.startsWith('/profile') && to.path.startsWith('/users/')) {
                this.loadProfile(to.params.username);
            }

            if (from.path.startsWith('/users/') && to.path.startsWith('/profile')) {
                this.pathUser = this.username;
                this.loadProfile(this.username);
            }

            if (to.path.startsWith('/users/') && to.params.username != this.pathUser) {
                this.pathUser = to.params.username;
                this.loadProfile(this.pathUser);
            }
        }
    },

	mounted() {
        if (this.pathUser){
            this.loadProfile(this.pathUser)
        } else {
            this.loadProfile(this.username)
        }
	}
}
</script>

<template>
    <!-- <h1 v-if="!this.logged" class="positive-banner" style="margin-top: 40px;">Please log-in to see profile contents</h1> -->
    <div>
        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
        <h2>Profile page</h2>
        <div class="d-flex align-items-center">
            <SearchBar />
			</div>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="loadProfile(pathUser)">
						Refresh
					</button>
				</div>
			</div>
		</div>
        
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        <div v-else class="profile-container">
            <div v-if="userNotFound" class="alert alert-warning" role="alert">
                User not found!
            </div>
            <LoadingSpinner :loading="loading" />
            <div class="info-photo-container">
                <div class="user-info">
                    <div class="user-name">
                        <h1>{{ pathUser || this.username }}</h1>
                    </div>
                    <div class="user-stats">
                        <p>Post: {{ numberOfPosts }}</p>
                        <p>Followers: {{ numberOfFollowers }}</p>
                        <p>Following: {{ numberOfFollowing }}</p>
                    </div>
                </div>

                <hr>

                <div v-if="this.profilePhotos" class="horizontal-photo-container">
                    <div v-for="photo in this.profilePhotos" :key="photo.PhotoID" class="horizontal-photo-div">

                        <PhotoCard 
                            :photo="photo"
                            @delete-post="deletePost(photo.PhotoID)"
                            @update-photo="updatePhoto"
                            />

                    </div>
                </div>
            </div>
            <div v-if="this.isOwner" class="user-actions-container">
                <div class="new-post">
                    <h3>New Post</h3>
                    <form id="post-form" @submit.prevent="newPost(description, picture)">
                        
                        Picture: <input type="file" v-on:change="fileUpload" /><br />
                        Description: <input type="text" v-model="description" /><br />
                        <br>
                        <div class="btn-group me-2">
                            <button type="button" class="btn btn-sm btn-outline-secondary" @click="newPost(description, picture)">
                                New Post
                            </button>
                        </div>
                    </form>
                </div>
                <hr>
                <div class="change-username">
                    <form id="username-form" @submit.prevent="changeUsername(newUsername)">
                        New Username: <input type="text" v-model="newUsername" /><br />
                        <button type="button" class="btn btn-sm btn-outline-secondary" @click="changeUsername(newUsername)">Change Username</button>
                    </form>
                    <div v-show="usernameTakenBanner" class="alert alert-warning" role="alert">
                        Username is already taken!
                    </div>
                </div>
            </div>
            <div v-else class="button-container">
                <div v-if="this.profileData.Followers">
                    <div v-if="followedByYou()">
                        <button id="unfollow-user-button" class="btn btn-sm btn-outline-primary" @click="unfollowUser()">Unfollow</button>
                    </div>
                    <div v-else>
                        <button id="follow-user-button" class="btn btn-sm btn-outline-primary" @click="followUser()">Follow</button>
                    </div>
                </div>
                <div v-if="this.profileData.Banned">
                    <button id="unfollow-user-button" class="btn btn-sm btn-outline-primary" @click="unbanUser()">Unban</button>
                </div>
                <div v-else>
                    <button id="follow-user-button" class="btn btn-sm btn-outline-primary" @click="banUser()">Ban</button>
                </div>
            </div>
        </div>
	</div>
</template>

<style>
.profile-container {
    display: flex;
    flex-direction: row;
}
.info-photo-container {
    width: 80%;
    margin-right: 30px;
}
/* maybe do it with a grid */
.horizontal-photo-container {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
}
.horizontal-photo-div {
    margin-right: 10px;
    margin-bottom: 20px;
}
.user-info {
    display: flex;
    align-items: center;
}
.user-name {
    flex: 1
}
.user-stats {
    flex: 4;
    align-items: center; 
    margin-left: 20px; 
}
.user-stats p {
    margin: 0; 
    margin-right: 10px; 
}
.user-actions-container {
    width: 30%;
}
.button-container {
  display: grid;
  grid-template-columns: repeat(2, 1fr); 
  grid-gap: 5px;
}
</style>
