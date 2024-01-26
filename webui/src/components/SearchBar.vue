<script>
import { RouterLink } from 'vue-router';

export default {

    components: {
        RouterLink,
    },

    emits: ['delete-post'],

    data: function() {
        return {
            searchUser: '',
        }
    },

    methods: {

        async handleSearch() {
			this.loading = true;
			this.errormsg = null;
            try {
                this.$router.push("/users/" + this.searchUser, { 
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
}
</script>

<template>
    <div class="search-bar">
        <input v-model="searchUser" @keyup.enter="handleSearch" placeholder="Search..." class="center-placeholder">
    </div>
</template>


<style scoped>
.search-bar input {
  padding: 8px;
  width: 500px;
  border: 1px solid #ccc;
  border-radius: 10px;
}
.center-placeholder::placeholder {
  text-align: center;
}
</style>