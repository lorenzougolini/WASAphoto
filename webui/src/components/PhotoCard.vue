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
        photoLocation: String,
        photoDescription: String,
        photoDate: String,
        photoLikes: Array,
        photoComments: Array,
        parent: String
    },

    methods: {
        formatDateFromUnix(unixDate) {
            let normalDate = new Date(unixDate * 1000);
            return normalDate.toLocaleString('it-EU');
        },

        buildUserLink() {
            return "/users/" + this.photoAuthor;
        },

        async toggleLike() {
            this.loading = true;   
            this.errormsg = null;
            
            try {         
                let response = await this.$axios.post("/photos/" + this.photoId + "/likes", {},{
                    headers: {
                        'Authorization': sessionStorage.getItem("userid"),
                    }
                });
                console.log(response);
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            console.log("Post liked!");
        },
    },

    computed: {
        photoLikesNum() {
            if (!this.photoLikes) {
                return 0;
            }
            return this.photoLikes.length;
        },
        photoCommentsNum() {
            if (!this.photoComments) {
                return 0;
            }
            return this.photoComments.length;
        }
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
                <button type="button" @click="toggleLike()"><svg class="feather"><use href="/feather-sprite-v4.29.0.svg#like"/></svg>Likes: {{ photoLikesNum }}</button>
                <button type="button"><svg class="feather"><use href="/feather-sprite-v4.29.0.svg#comment"/></svg>Comments: {{ photoCommentsNum }}</button>
                <!-- <p>Likes: {{ photoLikesNum }}</p>
                <p>Comments: {{ photoCommentsNum }}</p> -->
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