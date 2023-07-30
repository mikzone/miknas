<template>
  <q-list>
    <q-item>
      <q-item-section>
        <q-item-label caption>附件列表</q-item-label>
      </q-item-section>
      <q-item-section side>
        <q-btn dense round flat icon="add" @click.stop.prevent="">
          <q-menu>
            <q-list dense>
              <q-item
                v-for="attachConf, attachType in attachCfgs"
                :key="attachType"
                v-close-popup
                clickable
                @click="noteInfoUtil.startAddAttach(attachType)"
              >
                <q-item-section>
                  <q-item-label>{{ attachConf.name }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </q-item-section>
    </q-item>
    <component
      :is="attachCfgs[attachInfo.type].component"
      v-for="attachInfo in attachList"
      :key="attachInfo.id"
      :attach-info="attachInfo"
      :note-info-util = "noteInfoUtil"
    />
  </q-list>
</template>

<script setup>
import { computed } from 'vue';
import { attachCfgs, useNoteInfoUtil } from '../helpers';

const props = defineProps({
  noteInfo: {
    type: Object,
    required: true,
  },
});

const attachList = computed(() => {
  return props.noteInfo.noteAttachs || [];
});

const noteInfo2 = computed(()=>{
  return props.noteInfo;
})

const noteInfoUtil = useNoteInfoUtil(noteInfo2);

</script>
