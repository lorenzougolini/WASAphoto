<script>
import { RouterLink } from 'vue-router';

export default {

    components: {
        RouterLink,
    },

    emits: ['delete-post', 'update-photo'],

    props:{
        photo: {
            type: Object,
            required: true
        },
    },

    data: function() {
        return {
            loading: false,
            errormsg: null,
            showComments: false,
            comment: '',
            imageUrl: '',
        }
    },

    methods: {

        formatDateFromUnix(unixDate) {
            let normalDate = new Date(unixDate * 1000);
            return normalDate.toLocaleString('it-EU');
        },

        buildUserLink() {
            return "/users/" + this.photo.Author;
        },

        isLiked() {
            if (!this.photo.Likes) {
                return {"liked": false, "likeid": null};
            }

            for (let i = 0; i < this.photo.Likes.length; i++) {
                const like = this.photo.Likes[i]; 
                if (like.Username === sessionStorage.getItem("username")) {
                    return {"liked": true, "likeid": like.LikeID};
                }
            }
            return {"liked": false, "likeid": null};
        },

        toggleLike() {
            if (this.isLiked().liked) {
                this.unlikePhoto(this.isLiked().likeid);
            } else {
                this.likePhoto();
            }
        },

        isCommented() {
            if (!this.photo.Comments) {
                return false;
            }

            for (let i = 0; i < this.photo.Comments.length; i++) {
                const comment = this.photo.Comments[i]; 
                if (comment.Username === sessionStorage.getItem("username")) {
                    return true;
                }
            }
            return false;
        },

        toggleComments() {
            this.showComments = !this.showComments; 
        },

        canDeleteComment(comment) {
            return comment.Username === sessionStorage.getItem("username") || this.photo.Author === sessionStorage.getItem("username");
        },

        async likePhoto() {
            this.loading = true;   
            this.errormsg = null;

            try {         
                let response = await this.$axios.post("/photos/" + this.photo.PhotoID + "/likes", {},{
                    headers: {
                        'Authorization': sessionStorage.getItem("userid"),
                    }
                });

                this.$emit('update-photo', response.data);

                console.log("Post liked!");

            } catch (e) {
                this.errormsg = e.toString();

            }
            this.loading = false;
        },
        
        async unlikePhoto(likeid) {
            try {

                let response = await this.$axios.delete("/photos/" + this.photo.PhotoID + "/likes/" + likeid , {
                    headers: {
                        'Authorization': sessionStorage.getItem("userid"),
                    }
                });

                this.$emit('update-photo', response.data);

                console.log("Post disliked!");

            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
        },
        
        async commentPhoto() {
            this.loading = true;   
            this.errormsg = null;

            try {         
                let response = await this.$axios.post("/photos/" + this.photo.PhotoID + "/comments", {
                    "comment": this.comment
                },{
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        'Authorization': sessionStorage.getItem("userid"),
                    }
                });

                this.$emit('update-photo', response.data);

                this.comment = '';
                console.log("Comment posted!");

            } catch (e) {
                this.errormsg = e.toString();

            }
            this.loading = false;
        },
        
        async uncommentPhoto(commentid) {
            try {

                let response = await this.$axios.delete("/photos/" + this.photo.PhotoID + "/comments/" + commentid , {
                    headers: {
                        'Authorization': sessionStorage.getItem("userid"),
                    }
                });

                console.log(response.data);

                this.$emit('update-photo', response.data);

                console.log("Comment deleted!");

            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
        },

    },

    computed: {
        likesCount() {
            if (!this.photo.Likes) {
                return 0;
            }
            return this.photo.Likes.length;
        },
        commentCount() {
            if (!this.photo.Comments) {
                return 0;
            }
            return this.photo.Comments.length;
        },
        canDeletePhoto() {
            return this.photo.Author === sessionStorage.getItem("username");
        },
        decodedPhoto() {
            return "data:image/*;base64," + this.photo.File;
        },
    },

}
</script>
<template>

    <div class="photo-card">

        <div class="author-container">
            <div>
                <RouterLink :to="buildUserLink()" class="nav-link">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
                    {{ photo.Author }}
                </RouterLink>
                <p class="photo-date-rigth">{{ formatDateFromUnix(photo.DateAndTime) }}</p>
            </div>
        </div>

        <div class="image-container">
            <img :src="decodedPhoto" alt="Image not loaded" />
        </div>

        <div class="descdate-div">
            <br>
            <p>{{ photo.Description }}</p>
        </div>

        <div class="like-comment-div">

            <div class="btn-group">
                <button @click="toggleLike" :class="{ 'btn': true, 'btn-danger': this.isLiked().liked, 'btn-outline-danger': !this.isLiked().liked }">
                    Likes: {{ likesCount }}
                </button>
                <button @click="toggleComments" :class="{ 'btn': true, 'btn-primary': this.isCommented(), 'btn-outline-primary': !this.isCommented() }">
                    View comments: {{ commentCount }}
                </button>
            </div>

            <div v-show="showComments" class="comment-section">

                <div class="comments">
                    <div v-for="comment in photo.Comments" :key="comment.CommentID" class="comment">
                        <div class="comment-text">
                            <p><b>{{ comment.Username }}</b>: {{ comment.CommentText }}</p>
                            <hr/>
                        </div>
                        
                        <div v-show="canDeleteComment(comment)" class="delete-comment">
                            <button class="btn btn-sm" @click="uncommentPhoto(comment.CommentID)">
                                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#trash"/></svg>
                            </button>
                        </div>
                    </div>
                </div>
                <div>
                    <input id="comment-input" type="text" v-model="comment" />
                    <button class="btn btn-sm btn-outline-primary" @click="commentPhoto()">Comment</button>
                </div>
                
            </div>

        </div>

        <div v-show="canDeletePhoto" class="photo-delete">
            <button class="btn btn-sm btn-outline-danger" @click="$emit('delete-post')">Delete</button>
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
    height: 100%;
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

.comment-section {
    margin-top: 10px;
}
.comments{
    max-height: 150px; /* Set the maximum height for the scrollable area */
    overflow-y: auto; /* Enable vertical scrolling */
    border: 1px solid #ccc; /* Optional: Add a border for clarity */
    padding: 10px; /* Optional: Add padding for spacing */
}
.comment {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.comment-content {
    flex-grow: 1;
}

.delete-comment {
    margin-left: 10px;
}

.delete-comment button {
    padding: 0;
    cursor: pointer;
}

.author-container div {
    display: flex;
    align-items: center;
    padding: 3px;
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
.author-container div:hover {
    background-color: #e6e6e6;
}   
.image-container {
    height: 400px;
    overflow: hidden;
    display: flex;
    align-items: center;
    justify-content: center; 
    border-radius: 5px;
}

.image-container img {
    max-width: 100%; 
    max-height: 100%; 
    object-fit: contain;
}
.photo-delete {
    display: flex;
    justify-content: flex-end;
}

</style>