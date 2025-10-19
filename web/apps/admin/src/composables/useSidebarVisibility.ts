import { computed, ref } from 'vue';
import { preferences } from '@vben/preferences';

/**
 * Composable for managing sidebar visibility based on user preferences
 */
export function useSidebarVisibility() {
  // Get sidebar preferences
  const sidebarPrefs = computed(() => (preferences.sidebar as any) || {});
  
  // Get selected projects
  const selectedProjects = computed(() => sidebarPrefs.value.selectedProjects || []);
  
  // Get section visibility
  const sectionVisibility = computed(() => sidebarPrefs.value.sectionVisibility || {});
  
  // Check if a specific section should be visible
  const isSectionVisible = (sectionKey: string): boolean => {
    return sectionVisibility.value[sectionKey] ?? true;
  };
  
  // Check if a project should be visible
  const isProjectVisible = (projectId: string): boolean => {
    return selectedProjects.value.includes(projectId);
  };
  
  // Get visible sections
  const visibleSections = computed(() => {
    const sections = [
      { key: 'dashboard', label: 'Dashboard', icon: 'lucide:home' },
      { key: 'eln', label: 'ELN', icon: 'lucide:book-open' },
      { key: 'admin', label: 'Administration', icon: 'lucide:settings' },
      { key: 'projects', label: 'Projects', icon: 'lucide:folder' },
      { key: 'reports', label: 'Reports', icon: 'lucide:file-text' },
      { key: 'settings', label: 'Settings', icon: 'lucide:cog' },
    ];
    
    return sections.filter(section => isSectionVisible(section.key));
  });
  
  // Get visible projects
  const visibleProjects = computed(() => {
    // This would typically come from an API call
    // For now, return the selected projects
    return selectedProjects.value;
  });
  
  return {
    selectedProjects,
    sectionVisibility,
    isSectionVisible,
    isProjectVisible,
    visibleSections,
    visibleProjects,
  };
}
