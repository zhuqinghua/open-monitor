<template>
  <div class=" ">
    <div class="alarm-total">
      <Tag color="success"><span style="font-size:14px">Low:{{this.low}}</span></Tag>
      <Tag color="warning"><span style="font-size:14px">Medium:{{this.mid}}</span></Tag>
      <Tag color="error"><span style="font-size:14px">High:{{this.high}}</span></Tag>
    </div>
    <div class="alarm-list">
      <template v-for="(alarmItem, alarmIndex) in resultData">
        <section :key="alarmIndex" class="alarm-item c-dark-exclude-color" :class="'alarm-item-border-'+ alarmItem.s_priority">
          <i class="fa fa-bar-chart fa-operate" v-if="!alarmItem.is_custom" @click="goToEndpointView(alarmItem)" aria-hidden="true"></i>
          <ul>
            <li>
              <label class="alarm-item-label">{{$t('field.endpoint')}}:</label>
              <Tag type="border" color="primary">{{alarmItem.endpoint}}</Tag>
            </li>
            <li v-if="!alarmItem.is_custom">
              <label class="alarm-item-label">{{$t('field.metric')}}:</label>
              <Tag type="border" color="primary">{{alarmItem.s_metric}}</Tag>
            </li>
            <li>
              <label class="alarm-item-label">{{$t('tableKey.s_priority')}}:</label>
              <Tag type="border" color="primary">{{alarmItem.s_priority}}</Tag>
            </li>
            <li v-if="!alarmItem.is_custom && alarmItem.tags">
              <label class="alarm-item-label">Tags:</label>
              <Tag type="border" v-for="(t,tIndex) in alarmItem.tags.split('^')" :key="tIndex" color="cyan">{{t}}</Tag>
            </li>
            <li>
              <label class="alarm-item-label">{{$t('tableKey.start')}}:</label><span>{{alarmItem.start_string}}</span>
            </li>
            <li>
              <label class="alarm-item-label">{{$t('details')}}:</label>
              <span>
                <Tag color="default">{{$t('tableKey.start_value')}}:{{alarmItem.start_value}}</Tag>
                <Tag color="default" v-if="alarmItem.s_cond">{{$t('tableKey.threshold')}}:{{alarmItem.s_cond}}</Tag>
                <Tag color="default" v-if="alarmItem.s_last">{{$t('tableKey.s_last')}}:{{alarmItem.s_last}}</Tag>
                <Tag color="default" v-if="alarmItem.path">{{$t('tableKey.path')}}:{{alarmItem.path}}</Tag>
                <Tag color="default" v-if="alarmItem.keyword">{{$t('tableKey.keyword')}}:{{alarmItem.keyword}}</Tag>
              </span>
            </li>
            <li>
              <label class="alarm-item-label" style="vertical-align: top;">{{$t('alarmContent')}}:</label>
              <div class="col-md-9" style="display: inline-block;padding:0" v-html="alarmItem.content"></div>
            </li>
          </ul>
        </section>
      </template>
    </div>
  </div>
</template>

<script>
export default {
  name: '',
  data() {
    return {
      cutsomViewId: null,
      interval: '',

      resultData: [],
      low: 0,
      mid: 0,
      high: 0
    }
  },
  mounted () {},
  destroyed() {
    clearInterval(this.interval)
  },
  methods: {
    clearAlarmInterval () {
      clearInterval(this.interval)
    },
    getAlarm (id, viewCondition) {
      this.getAlarmdata(id)
      this.interval = setInterval(()=>{
        this.getAlarmdata(id)
      }, (viewCondition.autoRefresh || 10) * 1000)

    },
    getAlarmdata (id) {
      const parmas = {id}
      this.$root.$httpRequestEntrance.httpRequestEntrance('GET', '/monitor/api/v1/dashboard/custom/alarm/list', parmas, (responseData) => {
        this.resultData = responseData.data
        this.low = responseData.low
        this.mid = responseData.mid
        this.high = responseData.high
      }, {isNeedloading: false})
    },
    goToEndpointView (alarmItem) {
      const endpointObject = {
        option_value: alarmItem.endpoint,
        type: alarmItem.endpoint.split('_').slice(-1)[0]
      }
      this.$router.push({ name: 'endpointView',params: endpointObject})
    }
  },
  components: {},
}
</script>

<style scoped lang="less">
.alarm-item-label {
  width: 70px;
}
.alarm-empty {
  height: ~"calc(100vh - 180px)";
  width: ~"calc(100vw * 0.4)";
  text-align: center;
  padding:50px;
  color: #2d8cf0;
}
.flex-container {
  display: flex;
}
li {
  list-style: none;
} 

label {
  margin-bottom: 0;
  text-align: right;
}
.alarm-total {
  // float: right;
  margin-left: 8px;
  font-size: 18px;
}
.alarm-list {
  height: ~"calc(100vh - 220px)";
  width: 100%;
  overflow-y: auto;
}
.alarm-item {
  border: 1px solid @gray-d;
  margin: 8px;
  padding: 4px;
  border-radius: 4px;
  li {
    padding: 2px;
  }
}
.alarm-item-border-high {
  // border: 1px solid @color-orange-F;
  color: @color-orange-F;
}
.alarm-item-border-medium {
  // border: 1px solid @blue-2;
  color: @blue-2;
}
// .alarm-item-border-low {
//   // border: 1px solid @gray-d;
// }

.alarm-item:hover {
  box-shadow: 0 0 12px @gray-d;
}

.alarm-item /deep/.ivu-icon-ios-close:before {
  content: "\F102";
}

.fa-operate {
  margin: 8px;
  float: right;
  font-size: 16px;
  cursor: pointer;
}

/* 可以设置不同的进入和离开动画 */
/* 设置持续时间和动画函数 */
.slide-fade-enter-active {
  transition: all .3s ease;
}
.slide-fade-leave-active {
  transition: all .8s cubic-bezier(1.0, 0.5, 0.8, 1.0);
}
.slide-fade-enter, .slide-fade-leave-to
/* .slide-fade-leave-active for below version 2.1.8 */ {
  transform: translateX(10px);
  opacity: 0;
}
</style>
