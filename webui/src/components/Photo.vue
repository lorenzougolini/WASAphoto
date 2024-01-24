<script>
export default {
    emits: ['delete-post'],

    props:{
        photoId: String,
        photoLocation: String,
        photoDescription: String,
        photoDate: String,
        photoAuthor: String,
        photoLikes: Array,
        photoComments: Array
    },

    methods: {
        formatDateFromUnix(unixDate) {
            let normalDate = new Date(unixDate * 1000);
            return normalDate.toLocaleString();
        }
    },

    computed: {
        photoLikesNum() {
            if (photoLikes == null) {
                return 0;
            }
            return photoLikes.length;
        },
        photoCommentsNum() {
            if (photoComments == null) {
                return 0;
            }
            return photoComments.length;
        }
    }
}
</script>
<template>
    <div class="photo-card">
        <div class="author-container">
            <p>{{ photoAuthor }}</p>
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
                <button class="btn btn-sm btn-outline-danger" @click="$emit('delete-post')">Delete</button>
            </div>
        </div>
    </div>
</template>

<style>
.photo-card {
    border: 1px solid black;
    border-radius: 5px;
    padding: 5px;
    margin: 5px;
    width: 300px;
    /* height: 100%; */
    display: flex;
    flex-direction: column;
    justify-content: space-between;
}
.image-container {
    height: 300px;
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