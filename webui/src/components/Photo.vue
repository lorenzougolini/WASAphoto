<script>
export default {
    props:{
        photoLocation: String,
        photoDescription: String,
        photoDate: String,
    },

    methods: {
        format(date) {
            return new Date(date).toLocaleString();
        }
    }
}
</script>
<template>
    <div class="photo-card">
        <img :src="photoLocation" :alt='photoDescription'>
        <div class="photo-description">
            <p>{{ photoDescription }}</p>
        </div>
        <div class="photo-date">    
            <p>{{ format(photoDateAndTime) }}</p>
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
    height: 300px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
}
.photo-card img {
    width: 100%;
    height: 80%;
    object-fit: cover;
}
</style>