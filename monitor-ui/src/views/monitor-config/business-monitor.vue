<template>
  <div class=" ">
    <section>
      <ul class="search-ul">
        <li class="search-li">
          <Select
            style="width:300px"
            v-model="endpointID"
            filterable
            remote
            clearable
            :remote-method="getEndpointList"
            @on-change="changeEndpoint"
            @on-clear="clearEndpoint"
            >
            <Option v-for="(option, index) in endpointOptions" :value="option.id" :key="index">
             <Tag color="cyan" class="tag-width" v-if="option.type == 'host'">host</Tag>{{option.option_text}}</Option>
          </Select>
          <span @click="openDoc">
            <i 
              class="fa fa-book" 
              aria-hidden="true" 
              style="font-size:20px;color:#58a0e6;vertical-align: middle;margin-left:20px">
            </i>
            {{$t('operationDoc')}}
          </span>
        </li>
      </ul>
    </section>
    <section v-if="!!endpointID" style="margin-top: 16px;">
      <Tag color="blue">{{endpointGuid}}</Tag>
      <button @click="add" type="button" class="btn btn-sm btn-cancel-f">
        <i class="fa fa-plus"></i>
        {{$t('button.add')}}
      </button>
      <PageTable :pageConfig="pageConfig">
        <div slot='tableExtend'>
          <extendTable :detailConfig="pageConfig.table.isExtend.detailConfig"></extendTable>
        </div>
      </PageTable>
      <ModalComponent :modelConfig="modelConfig">
        <div slot="thresholdConfig" class="extentClass">  
          <div class="marginbottom params-each">
            <label class="col-md-2 label-name">{{$t('field.endpoint')}}:</label>
            <Select v-model="modelConfig.addRow.owner_endpoint" style="width:338px">
              <Option v-for="item in modelConfig.slotConfig.endpointOption" :value="item.guid" :key="item.guid">{{ item.guid }}</Option>
            </Select>
          </div>
        </div>
      </ModalComponent>
      <ModalComponent :modelConfig="ruleModelConfig">
        <div slot="ruleConfig" class="extentClass">
          <div class="marginbottom params-each">
           <div style="margin: 4px 12px;padding:8px 12px;border:1px solid #dcdee2;border-radius:4px">
              <template v-for="(item, index) in ruleModelConfig.addRow.metric_config">
                <p :key="index">
                  <Button
                    @click="deleterule('metric_config', index)"
                    size="small"
                    style="background-color: #ff9900;border-color: #ff9900;"
                    type="error"
                    icon="md-close"
                  ></Button>
                  <Input v-model="item.key" style="width: 146px" placeholder="key, e.g:[.*][.*]" />
                  <Input v-model="item.metric" style="width: 146px" placeholder="metric, e.g:code" />
                  <Select v-model="item.agg_type" filterable style="width:140px">
                    <Option v-for="agg in ruleModelConfig.slotConfig.aggOption" :value="agg" :key="agg">{{
                      agg
                    }}</Option>
                  </Select>
                </p>
              </template>
              <Button
                @click="addEmpty('metric_config')"
                type="success"
                size="small"
                style="background-color: #0080FF;border-color: #0080FF;"
                long
                >{{ $t('addMetricConfig') }}</Button
              >
            </div>
          </div>
          <div class="marginbottom params-each">
           <div style="margin: 4px 12px;padding:8px 12px;border:1px solid #dcdee2;border-radius:4px">
              <template v-for="(item, index) in ruleModelConfig.addRow.string_map">
                <p :key="index">
                  <Button
                    @click="deleterule('string_map', index)"
                    size="small"
                    style="background-color: #ff9900;border-color: #ff9900;"
                    type="error"
                    icon="md-close"
                  ></Button>
                  <Input v-model="item.key" style="width: 146px" placeholder="key" />
                  <Input v-model="item.string_value" style="width: 146px" placeholder="string_value" />
                  <InputNumber v-model="item.int_value" style="width: 140px"></InputNumber>
                </p>
              </template>
              <Button
                @click="addEmpty('string_map')"
                type="success"
                size="small"
                style="background-color: #0080FF;border-color: #0080FF;"
                long
                >{{ $t('addStringMap') }}</Button
              >
            </div>
          </div>
        </div>
      </ModalComponent>
      <Modal
        v-model="isShowWarning"
        title="Delete confirmation"
        @on-ok="ok"
        @on-cancel="cancel">
        <div class="modal-body" style="padding:30px">
          <div style="text-align:center">
            <p style="color: red">Will you delete it?</p>
          </div>
        </div>
      </Modal>
      <Modal
        v-model="isShowWarningDelete"
        title="Delete confirmation"
        @on-ok="okDelRow"
        @on-cancel="cancleDelRow">
        <div class="modal-body" style="padding:30px">
          <div style="text-align:center">
            <p style="color: red">Will you delete it?</p>
          </div>
        </div>
      </Modal>
    </section>
  </div>
</template>

<script>
import extendTable from '@/components/table-page/extend-table'
let tableEle = [
  {title: 'tableKey.path', value: 'path', display: true},
  {title: 'tableKey.endpoint', value: 'owner_endpoint', display: true}
]
const btn = [
  {btn_name: 'button.add', btn_func: 'singleAddF'},
  {btn_name: 'button.edit', btn_func: 'editF'},
  {btn_name: 'button.remove', btn_func: 'deleteConfirmModal'},
]

export default {
  name: '',
  data() {
    return {
      isShowWarning: false,
      requestParams: null,
      isShowWarningDelete: false,
      type: '',
      typeValue: 'endpoint',
      typeList: [
        {label: 'field.endpoint', value: 'endpoint'},
        {label: 'field.group', value: 'grp'}
      ],

      endpointID: null,
      endpointGuid: null,
      endpointOptions: [],

      totalPageConfig: [],
      pageConfig: {
        table: {
          tableData: [],
          tableEle: tableEle,
          // filterMoreBtn: 'filterMoreBtn',
          primaryKey: 'id',
          btn: btn,
          handleFloat:true,
          isExtend: {
            func: 'getExtendInfo',
            data: {},
            slot: 'tableExtend',
            detailConfig: [{
              isExtendF: true,
              title: '',
              config: [
                {title: 'tableKey.regular', value: 'regular', display: true},
                {title: 'tableKey.tags', value: 'tags', display: true},
                {title: 'table.action',btn:[
                  {btn_name: 'button.edit', btn_func: 'editRuleItem'},
                  {btn_name: 'button.remove', btn_func: 'delRuleconfirmModal'}
                ]}
              ],
              data: [1],
              scales: ['25%', '20%', '15%', '20%', '20%']
            }]
          }
        }
      },
      modelTip: {
        key: '',
        value: 'metric'
      },
      ruleModelConfig: {
        modalId: 'rule_Modal',
        isAdd: true,
        modalStyle: 'min-width:530px',
        modalTitle: 'rule',
        saveFunc: 'saveRule',
        config: [
          {label: 'tableKey.regular', value: 'regular', placeholder: 'tips.required', v_validate: 'required:true', disabled: false, type: 'text'},
          {label: 'tableKey.tags', value: 'tags', placeholder: '', disabled: false, type: 'text'},
          {name:'ruleConfig',type:'slot'}
        ],
        addRow: { // [通用]-保存用户新增、编辑时数据
          regular: null,
          tags: null,
          metric_config: [],
          string_map: []
        },
        slotConfig: {
          aggOption: ['sum', 'avg', 'count']
        }
      },
      modelConfig: {
        modalId: 'add_edit_Modal',
        modalTitle: 'field.businessMonitor',
        isAdd: true,
        config: [
          {label: 'tableKey.path', value: 'path', placeholder: 'tips.required', v_validate: 'required:true', disabled: false, type: 'text'},
          {name:'thresholdConfig',type:'slot'}
        ],
        addRow: { // [通用]-保存用户新增、编辑时数据
          path: null,
          owner_endpoint: ''
        },
        slotConfig: {
          endpointOption: []
        }
      },
      id: null,
      activeData: null,
      activeIndex: '',
      extendData: null,
    }
  },
  
  async mounted () {
    await this.getEndpointList('.')
    if (!this.$root.$validate.isEmpty_reset(this.$route.params)) {
       this.$parent.activeTab = '/monitorConfigIndex/businessMonitor'
      this.endpointID = Number(this.$route.params.id)
      this.endpointGuid = this.$route.params.guid
      this.requestData(this.endpointID)
    }
  },
  beforeDestroy () {
    this.$root.$store.commit('changeTableExtendActive', -1)
  },
  methods: {
    openDoc () {
      window.open('http://webankpartners.gitee.io/wecube-docs/manual-open-monitor-config/#_6')
    },
    /*********/
    addEmpty (type) {
      if (!this.ruleModelConfig.addRow[type]) {
        this.ruleModelConfig.addRow[type] = []
      }
      if (type === 'metric_config') {
        this.ruleModelConfig.addRow[type].push({
          key: '',
          metric: '',
          agg_type: 'avg'
        })
      } else {
        this.ruleModelConfig.addRow[type].push({
          key: '',
          string_value: '',
          int_value: 0
        })
      }
    },
    deleterule(type, index) {
      this.ruleModelConfig.addRow[type].splice(index, 1)
    },
    saveRule () {
      this.$validator.validate().then(result => {
        if (!result) return
        let params = null
        if (this.ruleModelConfig.isAdd === true) {
          let newRow = JSON.parse(JSON.stringify(this.ruleModelConfig.addRow))
          newRow.id = 0
          let newData = JSON.parse(JSON.stringify(this.activeData))
          newData.rules.push(newRow)
          const index = this.pageConfig.table.tableData.findIndex(item => item.id === this.activeData.id)
          this.pageConfig.table.tableData[index] = newData
          params = {
            endpoint_id: this.endpointID,
            path_list: this.pageConfig.table.tableData,
          }
        } else {
          // 获取上层数据序号
          const index = this.pageConfig.table.tableData.findIndex(item => item.id === this.activeData.pId)
          let copyData = JSON.parse(JSON.stringify(this.pageConfig.table.tableData))
          let allRules = copyData[index].rules
          const ruleIndex = allRules.findIndex(item => item.id === this.activeData.id)
          let newRow = JSON.parse(JSON.stringify(this.ruleModelConfig.addRow))
          allRules[ruleIndex] = newRow
          params = {
            endpoint_id: this.endpointID,
            path_list: copyData,
          }
        }
        this.$root.$httpRequestEntrance.httpRequestEntrance('POST', this.$root.apiCenter.businessMonitor.update, params, () => {
          this.$Message.success(this.$t('tips.success'))
          this.$root.JQ('#rule_Modal').modal('hide')
          this.requestData(this.endpointID)
        })
      })
    },
    /*********/
    typeChange() {
      this.totalPageConfig = []
      this.getEndpointList('.')
    },
    clearEndpoint () {
      this.totalPageConfig = []
      this.getEndpointList('.')
    },
    changeEndpoint (val) {
      if (val) {
        this.endpointGuid = this.endpointOptions.find(item => item.id === val).option_text
        this.requestData(this.endpointID)
      } else {
        this.endpointGuid = ''
      }
    },
    async getEndpointList (query) {
      const params = {type: 'endpoint', search: query}
      this.$root.$httpRequestEntrance.httpRequestEntrance('GET', this.$root.apiCenter.resourceSearch.strategyApi, params, (responseData) => {
        this.endpointOptions = responseData.filter(item => item.type === 'host')
      })
    },
    requestData (id) {
      let params = {id}
      this.totalPageConfig = []
      this.$root.$httpRequestEntrance.httpRequestEntrance('GET', this.$root.apiCenter.businessMonitor.list, params, (responseData) => {
        this.pageConfig.table.tableData = responseData.path_list
        this.$root.$store.commit('changeTableExtendActive', -1)
      })
    },
    deleteConfirmModal (rowData) {
      this.selectedData = rowData
      this.isShowWarning = true
    },
    ok () {
      this.delF(this.selectedData)
    },
    cancel () {
      this.isShowWarning = false
    },
    deleteConfirm (rowData) {
      this.$delConfirm({
        msg: rowData.path,
        callback: () => {
          this.delF(rowData)
        }
      })
    },
    delF (rowData) {
      const index = this.pageConfig.table.tableData.findIndex(item => item.id === rowData.id)
      let tableData = JSON.parse(JSON.stringify(this.pageConfig.table.tableData))
      tableData.splice(index,1)
      let params = {
        endpoint_id: this.endpointID,
        path_list: tableData
      }
      this.$root.$httpRequestEntrance.httpRequestEntrance('POST', this.$root.apiCenter.businessMonitor.update, params, () => {
        this.$Message.success(this.$t('tips.success'))
        this.requestData(this.endpointID)
      })
    },
    singleAddF (rowData) {
      this.activeData = rowData
      this.ruleModelConfig.isAdd = true
      this.$root.JQ('#rule_Modal').modal('show')
    },
    add () {
      this.modelConfig.isAdd = true
      this.$root.$httpRequestEntrance.httpRequestEntrance('GET', this.$root.apiCenter.endpointManagement.list.api + '?page=1&size=1000', '', responseData => {
        this.modelConfig.slotConfig.endpointOption = responseData.data
        this.$root.JQ('#add_edit_Modal').modal('show')
      })
    },
    addPost () {
      this.$validator.validate().then(result => {
        if (!result) return
        let params = {
          endpoint_id: this.endpointID,
          path_list: [
            {
              id: 0,
              ...this.modelConfig.addRow,
              rules: []
            }
          ]
        }
        this.$root.$httpRequestEntrance.httpRequestEntrance('POST', this.$root.apiCenter.businessMonitor.add, params, () => {
          this.$Message.success(this.$t('tips.success'))
          this.$root.JQ('#add_edit_Modal').modal('hide')
          this.requestData(this.endpointID)
        })
      })
    },
    editF (rowData) {
      this.modelConfig.isAdd = false
      this.activeData = rowData
      this.modelConfig.addRow.path = rowData.path
      this.modelConfig.addRow.owner_endpoint = rowData.owner_endpoint
      this.modelTip.value = rowData.path
      this.$root.$httpRequestEntrance.httpRequestEntrance('GET', this.$root.apiCenter.endpointManagement.list.api + '?page=1&size=1000', '', responseData => {
        this.modelConfig.slotConfig.endpointOption = responseData.data
        this.$root.JQ('#add_edit_Modal').modal('show')
      })
    },
    editPost () {
      this.$validator.validate().then(result => {
        if (!result) return
        const index = this.pageConfig.table.tableData.findIndex(item => item.id === this.activeData.id)
        let tableData = JSON.parse(JSON.stringify(this.pageConfig.table.tableData))
        tableData[index].path = this.modelConfig.addRow.path
        tableData[index].owner_endpoint = this.modelConfig.addRow.owner_endpoint
        let params = {
          endpoint_id: this.endpointID,
          path_list: tableData
        }
        this.$root.$httpRequestEntrance.httpRequestEntrance('POST', this.$root.apiCenter.businessMonitor.update, params, () => {
          this.$Message.success(this.$t('tips.success'))
          this.$root.JQ('#add_edit_Modal').modal('hide')
          this.requestData(this.endpointID)
        })
      })
    },
    getExtendInfo(item){
      item.rules.forEach(xx => xx.pId = item.id)
      this.pageConfig.table.isExtend.detailConfig[0].data = item.rules
    },
    editRuleItem (rowData) {
      this.ruleModelConfig.isAdd = false
      // 获取上层数据序号
      const index = this.pageConfig.table.tableData.findIndex(item => item.id === rowData.pId)
      let copyData = JSON.parse(JSON.stringify(this.pageConfig.table.tableData))
      let allRules = copyData[index].rules
      const ruleIndex = allRules.findIndex(item => item.id === rowData.id)
      this.ruleModelConfig.addRow = allRules[ruleIndex]
      this.activeData = allRules[ruleIndex]
      this.$root.JQ('#rule_Modal').modal('show')
    },
    delPathconfirm (rowData) {
      this.$delConfirm({
        msg: rowData.keyword,
        callback: () => {
          this.delPathItem(rowData)
        }
      })
    },
    delRuleconfirmModal (rowData) {
      this.selectedData = rowData
      this.isShowWarningDelete = true
    },
    okDelRow () {
      this.delRuleItem(this.selectedData)
    },
    cancleDelRow () {
      this.isShowWarningDelete = false
    },
    delRuleItem (rowData) {
      // 获取上层数据序号
      const index = this.pageConfig.table.tableData.findIndex(item => item.id === rowData.pId)
      let copyData = JSON.parse(JSON.stringify(this.pageConfig.table.tableData))
      let allRules = copyData[index].rules
      const ruleIndex = allRules.findIndex(item => item.id === rowData.id)
      allRules.splice(ruleIndex, 1)
      const params = {
        endpoint_id: this.endpointID,
        path_list: copyData,
      }
      this.$root.$httpRequestEntrance.httpRequestEntrance('POST', this.$root.apiCenter.businessMonitor.update, params, () => {
        this.$Message.success(this.$t('tips.success'))
        this.requestData(this.endpointID)
      })
    }
  },
  components: {
    extendTable
  },
}
</script>

<style scoped lang="less">
.search-li {
    display: inline-block;
  }
  .search-ul>li:not(:first-child) {
    padding-left: 10px;
  }
</style>
<style scoped lang="less">
  .is-danger {
    color: red;
    margin-bottom: 0px;
  }
  .search-input {
    height: 32px;
    padding: 4px 7px;
    font-size: 12px;
    border: 1px solid #dcdee2;
    border-radius: 4px;
    width: 230px;
  }

  .section-table-tip {
    margin: 24px 20px 0;
  }
  .search-input:focus {
    outline: 0;
    border-color: #57a3f3;
  }

  .search-input-content {
    display: inline-block;
    vertical-align: middle; 
  }
  .tag-width {
    cursor: auto;
    width: 55px;
    text-align: center;
  } 
</style>
