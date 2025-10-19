<template>
  <div class="sidebar-menu">
    <!-- Dashboard Section -->
    <div v-if="isSectionVisible('dashboard')" class="menu-section">
      <div class="menu-section-header">
        <h3 class="menu-section-title">Dashboard</h3>
      </div>
      <ul class="menu-items">
        <li class="menu-item">
          <router-link to="/dashboard/workspace" class="menu-link">
            <div class="menu-icon">
              <LucideHome class="w-5 h-5" />
            </div>
            <span class="menu-text">Workspace</span>
          </router-link>
        </li>
        <li class="menu-item">
          <router-link to="/dashboard/analytics" class="menu-link">
            <div class="menu-icon">
              <LucideBarChart3 class="w-5 h-5" />
            </div>
            <span class="menu-text">Analytics</span>
          </router-link>
        </li>
      </ul>
    </div>

    <!-- ELN Section -->
    <div v-if="isSectionVisible('eln')" class="menu-section">
      <div class="menu-section-header">
        <h3 class="menu-section-title">ELN</h3>
      </div>
      <ul class="menu-items">
        <li class="menu-item">
          <router-link to="/eln/newReaction" class="menu-link">
            <div class="menu-icon">
              <LucideFlaskConical class="w-5 h-5" />
            </div>
            <span class="menu-text">New Reaction</span>
          </router-link>
        </li>
        <li class="menu-item">
          <router-link to="/eln/select-my-project" class="menu-link">
            <div class="menu-icon">
              <LucideSettings class="w-5 h-5" />
            </div>
            <span class="menu-text">Project Selection</span>
          </router-link>
        </li>
      </ul>
    </div>

    <!-- Projects Section -->
    <div v-if="isSectionVisible('projects')" class="menu-section">
      <div class="menu-section-header">
        <h3 class="menu-section-title">Projects</h3>
      </div>
      <ul class="menu-items">
        <li v-for="project in visibleProjects" :key="project" class="menu-item">
          <router-link :to="`/projects/${project}`" class="menu-link">
            <div class="menu-icon">
              <LucideFolder class="w-5 h-5" />
            </div>
            <span class="menu-text">{{ project }}</span>
          </router-link>
        </li>
        <li v-if="visibleProjects.length === 0" class="menu-item">
          <div class="menu-link disabled">
            <div class="menu-icon">
              <LucideFolderOpen class="w-5 h-5" />
            </div>
            <span class="menu-text">No projects selected</span>
          </div>
        </li>
      </ul>
    </div>

    <!-- Administration Section -->
    <div v-if="isSectionVisible('admin')" class="menu-section">
      <div class="menu-section-header">
        <h3 class="menu-section-title">Administration</h3>
      </div>
      <ul class="menu-items">
        <li class="menu-item">
          <router-link to="/admin/users" class="menu-link">
            <div class="menu-icon">
              <LucideUsers class="w-5 h-5" />
            </div>
            <span class="menu-text">Users</span>
          </router-link>
        </li>
        <li class="menu-item">
          <router-link to="/admin/roles" class="menu-link">
            <div class="menu-icon">
              <LucideShield class="w-5 h-5" />
            </div>
            <span class="menu-text">Roles</span>
          </router-link>
        </li>
        <li class="menu-item">
          <router-link to="/admin/project" class="menu-link">
            <div class="menu-icon">
              <LucideFolderPlus class="w-5 h-5" />
            </div>
            <span class="menu-text">Project Management</span>
          </router-link>
        </li>
      </ul>
    </div>

    <!-- Reports Section -->
    <div v-if="isSectionVisible('reports')" class="menu-section">
      <div class="menu-section-header">
        <h3 class="menu-section-title">Reports</h3>
      </div>
      <ul class="menu-items">
        <li class="menu-item">
          <router-link to="/reports/experiments" class="menu-link">
            <div class="menu-icon">
              <LucideFileText class="w-5 h-5" />
            </div>
            <span class="menu-text">Experiment Reports</span>
          </router-link>
        </li>
        <li class="menu-item">
          <router-link to="/reports/analytics" class="menu-link">
            <div class="menu-icon">
              <LucideTrendingUp class="w-5 h-5" />
            </div>
            <span class="menu-text">Analytics Reports</span>
          </router-link>
        </li>
      </ul>
    </div>

    <!-- Settings Section -->
    <div v-if="isSectionVisible('settings')" class="menu-section">
      <div class="menu-section-header">
        <h3 class="menu-section-title">Settings</h3>
      </div>
      <ul class="menu-items">
        <li class="menu-item">
          <router-link to="/settings/preferences" class="menu-link">
            <div class="menu-icon">
              <LucideCog class="w-5 h-5" />
            </div>
            <span class="menu-text">Preferences</span>
          </router-link>
        </li>
        <li class="menu-item">
          <router-link to="/settings/profile" class="menu-link">
            <div class="menu-icon">
              <LucideUser class="w-5 h-5" />
            </div>
            <span class="menu-text">Profile</span>
          </router-link>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useSidebarVisibility } from '#/composables/useSidebarVisibility';
import {
  LucideHome,
  LucideBarChart3,
  LucideFlaskConical,
  LucideSettings,
  LucideFolder,
  LucideFolderOpen,
  LucideUsers,
  LucideShield,
  LucideFolderPlus,
  LucideFileText,
  LucideTrendingUp,
  LucideCog,
  LucideUser,
} from 'lucide-vue-next';

const { isSectionVisible, visibleProjects } = useSidebarVisibility();
</script>

<style scoped lang="scss">
.sidebar-menu {
  padding: 1rem;
  height: 100%;
  overflow-y: auto;
}

.menu-section {
  margin-bottom: 2rem;
  
  &:last-child {
    margin-bottom: 0;
  }
}

.menu-section-header {
  margin-bottom: 0.5rem;
}

.menu-section-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0;
}

.menu-items {
  list-style: none;
  padding: 0;
  margin: 0;
}

.menu-item {
  margin-bottom: 0.25rem;
  
  &:last-child {
    margin-bottom: 0;
  }
}

.menu-link {
  display: flex;
  align-items: center;
  padding: 0.75rem;
  border-radius: 0.5rem;
  text-decoration: none;
  color: #374151;
  transition: all 200ms ease-in-out;
  
  &:hover {
    background-color: #f3f4f6;
    color: #111827;
  }
  
  &.router-link-active {
    background-color: #dbeafe;
    color: #1d4ed8;
    font-weight: 500;
  }
  
  &.disabled {
    color: #9ca3af;
    cursor: not-allowed;
    
    &:hover {
      background-color: transparent;
    }
  }
}

.menu-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.25rem;
  height: 1.25rem;
  margin-right: 0.75rem;
  flex-shrink: 0;
}

.menu-text {
  font-size: 0.875rem;
  line-height: 1.25rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

// Custom scrollbar
.sidebar-menu {
  scrollbar-width: thin;
  scrollbar-color: rgba(0, 0, 0, 0.2) transparent;
  
  &::-webkit-scrollbar {
    width: 6px;
  }
  
  &::-webkit-scrollbar-track {
    background: transparent;
  }
  
  &::-webkit-scrollbar-thumb {
    background-color: rgba(0, 0, 0, 0.2);
    border-radius: 3px;
  }
  
  &::-webkit-scrollbar-thumb:hover {
    background-color: rgba(0, 0, 0, 0.3);
  }
}
</style>
