<script>
import { RouterLink } from 'vue-router';

export default {

    components: {
        RouterLink,
    },

    emits: ['delete-post'],

    props:{
        photoAuthor: String,
        photoId: String,
        photoDescription: String,
        photoDate: String,
        photoLikes: Array,
        photoComments: Array,
    },

    data: function() {
        return {
            loading: false,
            errormsg: null,
            like: this.checkLiked(),
        }
    },

    methods: {
        formatDateFromUnix(unixDate) {
            let normalDate = new Date(unixDate * 1000);
            return normalDate.toLocaleString('it-EU');
        },

        buildUserLink() {
            return "/users/" + this.photoAuthor;
        },

        checkLiked() {
            if (!this.photoLikes) {
                return {"liked": false, "likeid": null};
            }

            for (let i = 0; i < this.photoLikes.length; i++) {
                const like = this.photoLikes[i]; 
                if (like.Username === sessionStorage.getItem("username")) {
                    console.log("the likeid is: " + like.LikeID);
                    return {"liked": true, "likeid": like.LikeID};
                }
            }
            return {"liked": false, "likeid": null};
        },

        async toggleLike() {
            this.loading = true;   
            this.errormsg = null;

            if (!this.like.liked) {
                try {         
                    let response = await this.$axios.post("/photos/" + this.photoId + "/likes", {},{
                        headers: {
                            'Authorization': sessionStorage.getItem("userid"),
                        }
                    });

                    console.log(response);
                    console.log("Post liked!");

                } catch (e) {
                    this.errormsg = e.toString();

                }
                this.loading = false;

            } else {

                try {
                    console.log(this.photoId, this.like.likeid);

                    let response = await this.$axios.delete("/photos/" + this.photoId + "/likes/" + this.like.likeid , {
                        headers: {
                            'Authorization': sessionStorage.getItem("userid"),
                        }
                    });

                    // remove the like from photoLikes
                    this.photoLikes = this.photoLikes.filter(like => like.LikeID !== this.like.likeid);
                    console.log(this.photoLikes);

                    console.log(response);
                    console.log("Post disliked!");

                } catch (e) {
                    this.errormsg = e.toString();
                }
                this.loading = false;
            }
            
        },
    },

    computed: {
        likesCount() {
            if (!this.photoLikes) {
                return 0;
            }
            return this.photoLikes.length;
        },
        commentCount() {
            if (!this.photoComments) {
                return 0;
            }
            return this.photoComments.length;
        },
    }
}
</script>
<template>
    <div class="photo-card">
        <div class="author-container">
            <div>
                <RouterLink :to="buildUserLink()" class="nav-link">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
                    {{ photoAuthor }}
                </RouterLink>
                <p class="photo-date-rigth">{{ formatDateFromUnix(photoDate) }}</p>
            </div>
        </div>
        <div class="image-container">
            <img :src="photoLocation" :alt="photoDescription"><br>
        </div>
        <div class="descdate-div">
            <br>
            <p>{{ photoDescription }}</p>
        </div>
        <div class="like-comment-div">
            <div class="btn-group" role="group" aria-label="Basic example">
                <button type="button" @click="toggleLike()"><svg class="feather"><use href="/feather-sprite-v4.29.0.svg#like"/></svg>Likes: {{ likesCount }}</button>
                <button type="button"><svg class="feather"><use href="/feather-sprite-v4.29.0.svg#comment"/></svg>Comments: {{ commentCount }}</button>
            </div>
        </div>
        <div v-if="this.$route.params.username === photoAuthor" class="photo-delete">
            <button class="btn btn-sm btn-outline-danger" @search="$emit('delete-post')">Delete</button>
        </div>
    </div>
</template>

<style scoped>
.photo-card {
    border: 1px solid black;
    border-radius: 5px;
    padding: 5px;
    margin: 5px;
    width: 500px;
    /* height: 100%; */
    display: flex;
    flex-direction: column;
    justify-content: space-between;
}
.author-container {
    border: 1px solid gray;
    border-radius: 5px;
    display: flex;
    cursor: pointer;
}
.author-container div {
    display: flex;
    align-items: center;
    padding: 5px;
    font-size:16px
}
.author-container div .nav-link {
    margin-right: 250px;
}

.photo-date-rigth {
    font-size: 14px;
    font-family: sans-serif;
    margin: 0;
}
/* .author-container div:hover {
    background-color: #e6e6e6;
}    */
.image-container {
    height: 500px;
    overflow: hidden;
}
.image-container img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 5px;
}
.photo-delete {
    display: flex;
    justify-content: flex-end;
}
</style>