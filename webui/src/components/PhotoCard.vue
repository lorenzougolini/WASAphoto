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
        photoComments: Array
    },

    methods: {
        formatDateFromUnix(unixDate) {
            let normalDate = new Date(unixDate * 1000);
            return normalDate.toLocaleString('it-EU');
        },

        async loadProfile (photoAuthor) {
			this.loading = true;
			this.errormsg = null;
            try {
                this.$router.push("/users/" + photoAuthor, { 
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
            <div @click="loadProfile(photoAuthor)">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
                {{ photoAuthor }}
            </div>
        </div>
        <div class="image-container">
            <img :src="photoLocation" :alt="photoDescription"><br>
        </div>
        <div class="descdate-delete-div">
            <br>
            <p>{{ photoDescription }}</p>
            <p>{{ formatDateFromUnix(photoDate) }}</p>
            <p>Likes: {{ photoLikesNum }}</p>
            <p>Comments: {{ photoCommentsNum }}</p>
            <div class="photo-delete">
                <!-- <button class="btn btn-sm btn-outline-danger" @click="this.$parent.deletePost(photoId)">Delete</button> -->
                <button class="btn btn-sm btn-outline-danger" @search="$emit('delete-post')">Delete</button>
            </div>
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