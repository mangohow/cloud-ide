
<template>
    <div class="card">
      <div class="img-area">
        <img class="avatar" :src="info.avatar" alt="">
      </div>
      <div class="info-area">
        <div class="name">{{info.name}}</div>
        <div class="desc">{{info.desc}}</div>
        <div class="tags">
          <span v-for="(tag, i) in info.tags" :key="i">{{tag}}</span>
        </div>
      </div>
    </div>
    

</template>



<script>
export default {
    props: [
        "info"
    ],
    data() {
        return {
            
        }
    },
    methods: {

    }
}
</script>




<style lang="less" scoped>

.card {
  background-color: rgb(53, 57, 65);
  border-radius: 5px;
  height: 80px;
  padding: 24px 0;
  cursor: pointer;
}

.card:hover {
  transform:  scale(1.02);
  box-shadow: 0 14px 28px rgba(0, 0, 0, 0.25), 0 10px 10px rgba(0, 0, 0, 0.22);
  border-left: 10px pink;
}

.img-area {
  width: 80px;
  height: 100%;
  float: left;

  .avatar {
    width: 44px;
    height: 44px;
    margin: 18px;
  }
}

.info-area {
  float: left;
  width: calc(100% - 80px);
  height: 100%;
  text-align: left;

  .name {
    color: white;
    font-size: 20px;
    line-height: 30px;
  }

  .desc {
    color: rgb(154, 156, 161);
    font-size: 12px;
    line-height: 24px;
  }

  .tags {
    color: white;
    font-size: 12px;
    display: block;
    padding: 15px 0 10px 0;

    span {
      background-color: rgb(112, 116, 122);
      padding: 2px 8px;
      border-radius: 3px;
      margin-right: 5px;
      margin-left: 5px;
      text-align: center;
    }
  }
}


</style>