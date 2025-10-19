# Project Selection & Sidebar Configuration

This component provides a comprehensive interface for users to select which projects to display in the sidebar and configure sidebar section visibility.

## Features

### 🎯 **Project Selection**
- **Search & Filter**: Real-time search through available projects
- **Bulk Actions**: Select all or clear all projects with one click
- **Visual Feedback**: Clear indication of selected projects with highlighting
- **Project Details**: Display project name, description, status, and creator
- **Selection Summary**: Shows count of selected projects

### 🎛️ **Sidebar Configuration**
- **Section Toggle**: Enable/disable entire sidebar sections
- **Visual Indicators**: Eye icons show visibility status
- **Default States**: Sensible defaults for each section
- **Real-time Updates**: Changes apply immediately

### 💾 **Persistence**
- **User Preferences**: Settings saved to vben-admin preferences system
- **Cross-session**: Settings persist across browser sessions
- **Reset Functionality**: Easy reset to default configuration

## Usage

### Basic Integration

```vue
<template>
  <div>
    <!-- Your existing layout -->
    <SidebarMenu />
  </div>
</template>

<script setup>
import SidebarMenu from '#/components/SidebarMenu.vue';
</script>
```

### Using the Composable

```vue
<script setup>
import { useSidebarVisibility } from '#/composables/useSidebarVisibility';

const { 
  isSectionVisible, 
  visibleProjects, 
  selectedProjects 
} = useSidebarVisibility();

// Check if a section should be visible
const showDashboard = isSectionVisible('dashboard');

// Get list of visible projects
const projects = visibleProjects.value;
</script>
```

## Sidebar Sections

| Section | Key | Default | Description |
|---------|-----|---------|-------------|
| Dashboard | `dashboard` | ✅ | Analytics and overview |
| ELN | `eln` | ✅ | Lab experiments and reactions |
| Administration | `admin` | ❌ | System management and users |
| Projects | `projects` | ✅ | Project management |
| Reports | `reports` | ❌ | Data analysis and reporting |
| Settings | `settings` | ❌ | Application preferences |

## API Integration

The component integrates with the existing vben-admin API:

```typescript
// Load projects
const response = await getProjectList({
  page: 1,
  pageSize: 100,
});

// Save preferences
await updatePreferences({
  sidebar: {
    selectedProjects: ['project1', 'project2'],
    sectionVisibility: {
      dashboard: true,
      eln: true,
      admin: false,
    }
  }
});
```

## Styling

The component uses vben-admin design patterns:

- **Tailwind CSS**: For responsive layout and spacing
- **Element Plus**: For form components and interactions
- **Lucide Icons**: For consistent iconography
- **Custom SCSS**: For hover effects and transitions

## Responsive Design

- **Mobile**: Single column layout with stacked cards
- **Tablet**: Two-column layout with side-by-side cards
- **Desktop**: Full layout with optimal spacing

## Accessibility

- **Keyboard Navigation**: Full keyboard support
- **Screen Readers**: Proper ARIA labels and descriptions
- **Focus Management**: Clear focus indicators
- **Color Contrast**: Meets WCAG guidelines

## Future Enhancements

- [ ] Drag & drop project reordering
- [ ] Custom section creation
- [ ] Project grouping and categories
- [ ] Advanced filtering options
- [ ] Export/import configurations
