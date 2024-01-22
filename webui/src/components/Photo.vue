<script>
export default {
    props:{
        photo: Object,
    }
}
</script>
<template>
    <div class="photo-card">
        <img :src="photoLocation" :alt="photo.Description">
        <div class="photo-description">
            <p>{{ photo.Description }}</p>
        </div>
        <div class="photo-date">    
            <p>{{ format(photo.DateAndTime) }}</p>
        </div>
    </div>
</template>