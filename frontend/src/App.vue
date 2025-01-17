<script lang="ts" setup>
import dayjs from 'dayjs'
import timezone from 'dayjs/plugin/timezone'
import utc from 'dayjs/plugin/utc'
import { reactive, onMounted, computed, ref } from 'vue'
import { EventsOn } from '@wails/runtime'
import { GetConfig } from '@wails/go/main/App'
import JsonTools from '@/components/pages/JsonTools.vue'
import Beam from '@/components/pages/Beam.vue'
import Tinker from '@/components/pages/Tinker.vue'
import ProjectSettings from '@/components/pages/ProjectSettings.vue'
import AppHeader from '@/components/AppHeader.vue'
import AppTab from '@/components/AppTab.vue'
import AppTabBar from '@/components/AppTabBar.vue'
import NotificationContainer from '@/components/NotificationContainer.vue'
import LineWrapIcon from '@/components/icons/LineWrapIcon.vue'
import './style.css'

dayjs.extend(utc)
dayjs.extend(timezone)
dayjs.locale('pl')
dayjs.tz.setDefault("Europe/Warsaw")

const tabs = [
  { id: "tinker", title: 'Tinker', content: Tinker },
  { id: "jsontools", title: 'JSON Tools', content: JsonTools },
  { id: "beam", title: 'Beam', content: Beam },
]

const activeTab = ref(tabs[0].id)

const lineWrap = ref(false)

const enum Section {
  App = 'app',
  ProjectManager = 'project-manager',
}

const data = reactive({
  currentSection: Section.App,
  tabs: tabs,
  value: tabs[0].title,
  messages: [] as Array<BeamMessage>,
  appConfig: {} as AppConfig,
  showNotification: false,
})

EventsOn("beamMessage", function (messageData: RawBeamMessage) {
  data.messages.unshift({
    timestamp: dayjs().toDate().toLocaleTimeString(),
    payload: messageData.Payload,
    group: messageData.Group,
  })
})

EventsOn("configUpdated", function (appConfig: AppConfig) {
  data.appConfig = appConfig
})

const clearMessages = () => data.messages = []

const currentProject = computed<Project>(function () {
  return data.appConfig.Projects
    ? data.appConfig.Projects[data.appConfig.Currentproject]
    : {} as Project
})

const currentBadgeColor = computed<string>(function() {
  return getBadgeColor(currentProject.value?.Tag)
})

const getBadgeColor = (tag: string): string => {
  const foundTag = data.appConfig.Tags?.find(function(tagConfig) {
    if (tag === tagConfig.Label) {
      return tagConfig.Color
    }
  })

  return foundTag ? foundTag.Color : "black"
}

const openProjectManager = () => {
  data.currentSection = data.currentSection === Section.ProjectManager ? Section.App : Section.ProjectManager
}

const refreshAppConfig = () => {
  GetConfig().then((config) => {
    data.appConfig = config
  })
}

onMounted(() => refreshAppConfig())
</script>

<template>
  <div class="h-full overflow-auto">
    <NotificationContainer />

    <div class="h-full w-full flex flex-col overflow-hidden">
      <AppHeader 
        :current-project="currentProject" 
        :open-project-manager="openProjectManager" 
        :badge-color="currentBadgeColor"
        />

      <div v-if="data.currentSection !== Section.ProjectManager" class="flex px-2 border-b-2 border-gray-300 bg-gray-200 justify-between">
        <AppTabBar>
          <AppTab v-for="tab in tabs" 
            :active="activeTab == tab.id" 
            @click="activeTab = tab.id"
          >
            {{ tab.title }}
          </AppTab>
        </AppTabBar>
        <div class="grid items-center">
          <div class="text-gray-500 mr-2 select-none">
            <label>
              <input type="checkbox" class="hidden" v-model="lineWrap" />
              <LineWrapIcon :class="['w-6 cursor-pointer', lineWrap ? 'text-black' : 'text-gray-400']" title="toggle line wrap" />
            </label>
          </div>
        </div>
      </div>
      <div v-if="data.currentSection === Section.App" 
        class="h-full"
      >
        <Beam v-if="activeTab === 'beam'" :project="currentProject" :messages="data.messages" @clear-beam-messages="clearMessages" />
        <Tinker v-else-if="activeTab === 'tinker'" :project="currentProject" :lineWrap="lineWrap" />
        <JsonTools v-else-if="activeTab === 'jsontools'" :project="currentProject" :lineWrap="lineWrap" />
      </div>

      <div v-else-if="data.currentSection === Section.ProjectManager" 
        class="pb-5 h-full m-3"
      >
        <ProjectSettings
          :app-config="data.appConfig"
          :current-project="currentProject"
          @close="data.currentSection = Section.App"
          />
      </div>
    </div>
  </div>
</template>
