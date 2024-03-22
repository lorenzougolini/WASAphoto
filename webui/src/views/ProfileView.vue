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
            shownUser: this.$route.params.username,
            profileJson: {},
            banned: false,
            description: '',
            newUsername: '',    
		}
	},

	methods: {
		async loadProfile(shownUser) {
			this.loading = true;
			this.errormsg = null;
            try {

                let response = await this.$axios.get("/users/" + shownUser, {
                    headers: {
                        'Authorization': this.userid,
                    }
                });
                this.profileJson = response.data;
                this.banned = this.profileJson.Banned;

            } catch (e) {
                this.errormsg = e.toString();
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
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            console.log("Post deleted!");
            this.loadProfile(this.username);
        },

        async changeUsername(newUsername) {
            this.loading = true;   
            this.errormsg = null;
            
            try {        
                let response = await this.$axios.put("/users/" + this.username, newUsername, {
                    headers: {
                        'Authorization': this.userid,
                    }
                });

                console.log(response.data.Username);
                sessionStorage.setItem("username", response.data.Username);
                
                console.log("Username changed!");
                this.shownUser = response.data.Username;
                this.username = response.data.Username;
                this.$router.push("/users/" + response.data.Username);
                this.loadProfile(this.shownUser);

            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
        },

        async followUser(){
            this.loading = true;   
            this.errormsg = null;
            
            try {        
                await this.$axios.put("/users/" + this.username + "/followers/" + this.shownUser, null, {
                    headers: {
                        'Authorization': this.userid,
                    }
                });
                console.log("User followed!");
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            this.loadProfile(this.shownUser);
        },

        async unfollowUser() {
            this.loading = true;   
            this.errormsg = null;
            
            try {        
                await this.$axios.delete("/users/" + this.username + "/followers/" + this.shownUser, {
                    headers: {
                        'Authorization': this.userid,
                    }
                });
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            console.log("User unfollowed!");
            this.loadProfile(this.shownUser);            
        },

        async banUser() {
            this.loading = true;   
            this.errormsg = null;
            
            try {        
                await this.$axios.put("/users/" + this.username + "/banned/" + this.shownUser, null, {
                    headers: {
                        'Authorization': this.userid,
                    }
                });
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            console.log("User banned!");
            this.loadProfile(this.shownUser);
        },
        
        async unbanUser() {
            this.loading = true;   
            this.errormsg = null;
            
            try {        
                await this.$axios.delete("/users/" + this.username + "/banned/" + this.shownUser, {
                    headers: {
                        'Authorization': this.userid,
                    }
                });
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            console.log("User unbanned!");
            this.loadProfile(this.shownUser);
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
		this.loadProfile(this.shownUser)
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h2>Profile page</h2>
            <div class="d-flex align-items-center">
				<SearchBar @search-user="loadProfile(this.shownUser)"/>
			</div>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="loadProfile(shownUser)">
						Refresh
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        <div class="profile-container">
            <LoadingSpinner :loading="loading" />
            <div class="info-photo-container">
                <div class="user-info">
                    <div class="user-name">
                        <h1>{{ shownUser }}</h1>
                    </div>
                    <div class="user-stats">
                        <p>Post: {{ numberOfPosts }}</p>
                        <p>Followers: {{ numberOfFollowers }}</p>
                        <p>Following: {{ numberOfFollowing }}</p>
                    </div>
                </div>

                <hr>

                <div v-if="this.profileJson.Posts" class="horizontal-photo-container">
                    <div v-for="photo in this.profilePhotos" :key="photo.PhotoID" class="horizontal-photo-div">

                        <PhotoCard @delete-post="deletePost(photo.PhotoID)"
                            :photo="photo"/>

                    </div>
                </div>
            </div>
            <div class="vertical-line"></div>
            <div v-if="shownUser == this.username" class="user-actions-container">
                <div class="new-post">
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
                <hr>
                <div class="change-username">
                    <form>
                        New Username: <input type="text" v-model="newUsername" /><br />
                        <button type="button" @click="changeUsername(newUsername)">Change Username</button>
                    </form>
                </div>
            </div>
            <div v-else class="button-container">
                <div v-if="this.profileJson.Followers">
                    <div v-if="followedByYou()">
                        <button id="unfollow-user-button" class="btn btn-sm btn-outline-primary" @click="unfollowUser()">Unfollow</button>
                    </div>
                    <div v-else>
                        <button id="follow-user-button" class="btn btn-sm btn-outline-primary" @click="followUser()">Follow</button>
                    </div>
                </div>
                <div v-if="this.banned">
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
