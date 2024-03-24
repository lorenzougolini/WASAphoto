<script>

export default {

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

                this.searchUser = '';
                
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
		},
    },
}
</script>

<template>
    <div class="search-container">
        <div class="search-bar">
            <input v-model="searchUser" @keyup.enter="handleSearch" placeholder="Search user..." class="center-placeholder">
            <button @click="handleSearch" class="btn">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
            </button>
        </div>
    </div>
</template>


<style scoped>

.search-container {
  display: flex;
  border: 1px solid #ccc;
  border-radius: 4px;
}
.search-bar input {
  padding: 8px;
  width: 500px;
  border: none;
  outline: none;
  border-radius: 10px;
}
.center-placeholder::placeholder {
  text-align: left;
}
.search-bar button {
  padding: 8px;
  border: none;
  border-radius: 0 4px 4px 0;
}
</style>