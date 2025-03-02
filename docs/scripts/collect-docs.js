import { readFile, writeFile, mkdir, copyFile as fsCopyFile, rm } from 'fs/promises';
import { dirname, join, basename, extname, sep } from 'path';
import { fileURLToPath } from 'url';
import { glob } from 'glob';
import { watch } from 'fs';
import { exec } from 'child_process';
import { promisify } from 'util';
import sharp from 'sharp';
import os from 'os';

// Check if we're running on Windows
const isWindows = os.platform() === 'win32';

// Helper function to normalize paths across platforms
function normalizePath(path) {
  // Convert backslashes to forward slashes for consistency
  return path.replace(/\\/g, '/');
}

const execAsync = promisify(exec);
const __dirname = dirname(fileURLToPath(import.meta.url));
const rootDir = join(__dirname, '../../');
const docsDir = join(__dirname, '../src/content/docs');
const assetsDir = join(__dirname, '../src/content/docs/assets');

// Patterns to ignore
const ignorePatterns = [
  '**/node_modules/**',
  '**/docs/**',
  '**/dist/**',
  '**/vendor/**',
  '**/build/**',
  '**/tmp/**',
  '**/connectors/**', // Ignore connectors directory for general copying
  '**/plugins/**', // Ignore plugins directory for general copying
];

function generateTitle(filePath) {
  // Get the name of the directory containing README
  const dir = dirname(filePath);
  if (dir === '.') return 'Root Documentation';

  // Split the path into parts and take the last directory
  // Normalize the path first to ensure consistent separator usage
  const normalizedPath = normalizePath(dir);
  const parts = normalizedPath.split('/');
  const lastDir = parts[parts.length - 1];

  // Transform kebab-case or snake_case to Title Case
  return lastDir
    .replace(/[-_]/g, ' ')
    .replace(/([a-z])([A-Z])/g, '$1 $2') // Split camelCase
    .split(' ')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ');
}

function generateDescription(filePath) {
  // Create description based on file path
  // Normalize the path first
  const normalizedPath = normalizePath(filePath);
  const parts = normalizedPath.split('/').filter((part) => part !== 'README.md');
  if (parts.length === 0) return 'Main project documentation';

  return `Documentation for the ${parts.join(' ')} component`;
}

function addFrontMatter(content, filePath) {
  // If content already has frontmatter, parse it
  if (content.startsWith('---')) {
    return content;
  }

  const title = generateTitle(filePath);
  const description = generateDescription(filePath);

  const frontMatter = `---
title: ${title}
description: ${description}
---

`;

  return frontMatter + content;
}

function simplifyPath(filePath) {
  if (basename(filePath) === 'index.md') {
    return join(dirname(filePath), '..') + '.md';
  }
  return filePath;
}

async function copyFile(file) {
  // Normalize file path first
  const normalizedFile = normalizePath(file);
  const sourcePath = join(rootDir, normalizedFile);
  let targetDir = join(docsDir, dirname(normalizedFile));
  let targetPath = join(targetDir, 'index.md');
  
  // Normalize path for parts checking
  const normalizedTargetPath = normalizePath(targetPath);
  const parts = normalizedTargetPath.split('/');

  if (normalizedTargetPath && (parts.includes('connectors') || parts.includes('plugins'))) {
    const parentDir = parts[parts.length - 2];

    if (parentDir !== 'connectors' && parentDir !== 'plugins') {
      targetPath = join(dirname(targetPath), '.') + '.md';
    }
  }

  // Special handling for connectors and plugins
  // Use normalized paths for includes checks
  if ((normalizedFile.includes('/connectors/') || normalizedFile.includes('/plugins/')) && 
      normalizedFile.toLowerCase().endsWith('readme.md')) {
    const pathParts = normalizedFile.split('/');
    const typeIndex = pathParts.findIndex((part) => part === 'connectors' || part === 'plugins');
    if (typeIndex !== -1 && typeIndex + 1 < pathParts.length) {
      const type = pathParts[typeIndex];
      const name = pathParts[typeIndex + 1];
      targetDir = join(docsDir, type);
      targetPath = join(targetDir, `${name}.md`);
    }
  }

  try {
    // Create target directory
    await mkdir(targetDir, { recursive: true });

    // Read the file
    const content = await readFile(sourcePath, 'utf8');

    // Add frontmatter if needed
    const processedContent = addFrontMatter(content, normalizedFile);

    // Write the processed content
    await writeFile(targetPath, processedContent);

    console.log(`Copied and processed ${normalizedFile} to ${targetPath}`);
  } catch (error) {
    console.error(`Error copying file ${normalizedFile}:`, error);
  }
}

async function processAndCopyImage(sourcePath, targetPath) {
  try {
    const ext = extname(sourcePath).toLowerCase();
    const isImage = ['.png', '.jpg', '.jpeg', '.gif', '.webp'].includes(ext);

    if (!isImage) {
      // If not an image, just copy the file
      await fsCopyFile(sourcePath, targetPath);
      return;
    }

    // For GIF files, convert to WebP while preserving animation
    if (ext === '.gif') {
      const targetWebP = targetPath.replace(/\.gif$/i, '.webp');
      await sharp(sourcePath, { animated: true }).webp({ quality: 80, effort: 6 }).toFile(targetWebP);
      console.log(`Converted GIF to WebP: ${sourcePath} -> ${targetWebP}`);
      return;
    }

    // Process other image types
    await sharp(sourcePath)
      .resize(1200, 1200, {
        // Maximum dimensions
        fit: 'inside', // Preserve aspect ratio
        withoutEnlargement: true, // Don't enlarge small images
      })
      .toFile(targetPath);

    console.log(`Processed and copied image ${sourcePath} to ${targetPath}`);
  } catch (error) {
    console.error(`Error processing image ${sourcePath}:`, error);
    // If processing failed, try to just copy the file
    await fsCopyFile(sourcePath, targetPath);
  }
}

async function copyAssets() {
  try {
    // Find all files in assets directory
    const assetFiles = await glob('**/assets/**/*.*', {
      cwd: rootDir,
      ignore: ignorePatterns,
      nocase: true,
    });

    const assetFilesContent = await glob('**/docs/content/**/*.md', {
      cwd: rootDir,
      nocase: true,
    });

    console.log('Found asset files:', assetFiles);
    console.log('Found asset files content:', assetFilesContent);

    for (const file of assetFiles) {
      const normalizedFile = normalizePath(file);
      const sourcePath = join(rootDir, normalizedFile);
      const targetPath = join(assetsDir, basename(normalizedFile));

      try {
        // Create assets directory if it doesn't exist
        await mkdir(dirname(targetPath), { recursive: true });

        // Process and copy the asset file
        await processAndCopyImage(sourcePath, targetPath);
      } catch (error) {
        console.error(`Error copying asset ${normalizedFile}:`, error);
      }
    }

    for (const file of assetFilesContent) {
      const normalizedFile = normalizePath(file);
      const sourcePath = join(rootDir, normalizedFile);
      const targetPath = join(docsDir, normalizedFile.replace('docs/', ''));

      try {
        // Create assets directory if it doesn't exist
        await mkdir(dirname(targetPath), { recursive: true });

        // Process and copy the asset file
        await fsCopyFile(sourcePath, targetPath);
      } catch (error) {
        console.error(`Error copying asset ${normalizedFile}:`, error);
      }
    }
  } catch (error) {
    console.error('Error collecting assets:', error);
  }
}

async function collectConnectorsDocs() {
  try {
    const connectorsPath = join(rootDir, 'connectors');
    const connectorDirs = await glob('*/', {
      cwd: connectorsPath,
      nocase: true,
    });

    console.log('Found connector directories:', connectorDirs);

    // Create directory for connectors
    const connectorsDocsDir = join(docsDir, 'connectors');
    await mkdir(connectorsDocsDir, { recursive: true });

    // Process each connector
    for (const connectorDir of connectorDirs) {
      const connectorName = normalizePath(connectorDir).replace(/\/$/, ''); // Remove trailing slash
      const readmeFiles = await glob('readme.md', {
        cwd: join(connectorsPath, connectorDir),
        nocase: true,
      });
      const readmePath = readmeFiles.length > 0 ? join(connectorsPath, connectorDir, readmeFiles[0]) : join(connectorsPath, connectorDir, 'README.md');

      try {
        const content = await readFile(readmePath, 'utf8');
        const connectorDocPath = join(connectorsDocsDir, `${connectorName}.md`);

        // Add frontmatter and write content
        // Ensure path is normalized for frontmatter generation
        await writeFile(connectorDocPath, addFrontMatter(content, normalizePath(`connectors/${connectorName}/README.md`)));
        console.log(`Generated documentation for connector ${connectorName}`);
      } catch (error) {
        if (error.code === 'ENOENT') {
          console.log(`No README.md found for connector ${connectorName}`);
        } else {
          console.error(`Error processing connector ${connectorName}:`, error);
        }
      }
    }

    // Create index file for connectors
    const indexContent = `---
title: Connectors
description: List of all available connectors and their documentation
---

# Available Connectors

${connectorDirs
  .map((dir) => {
    const name = normalizePath(dir).replace(/\/$/, '');
    return `- [${name}](${name})`;
  })
  .join('\n')}
`;

    await writeFile(join(connectorsDocsDir, 'index.md'), indexContent);
    console.log('Generated connectors index');
  } catch (error) {
    console.error('Error collecting connectors documentation:', error);
  }
}

async function collectPluginsDocs() {
  try {
    const pluginsPath = join(rootDir, 'plugins');
    const pluginDirs = await glob('*/', {
      cwd: pluginsPath,
      nocase: true,
    });

    console.log('Found plugin directories:', pluginDirs);

    // Create directory for plugins
    const pluginsDocsDir = join(docsDir, 'plugins');
    await mkdir(pluginsDocsDir, { recursive: true });

    // Process each plugin
    for (const pluginDir of pluginDirs) {
      const pluginName = normalizePath(pluginDir).replace(/\/$/, ''); // Remove trailing slash
      const readmePath = join(pluginsPath, pluginDir, 'README.md');

      try {
        const content = await readFile(readmePath, 'utf8');
        const pluginDocPath = join(pluginsDocsDir, `${pluginName}.md`);

        // Add frontmatter and write content
        await writeFile(pluginDocPath, addFrontMatter(content, normalizePath(`plugins/${pluginName}/README.md`)));
        console.log(`Generated documentation for plugin ${pluginName}`);
      } catch (error) {
        if (error.code === 'ENOENT') {
          console.log(`No README.md found for plugin ${pluginName}`);
        } else {
          console.error(`Error processing plugin ${pluginName}:`, error);
        }
      }
    }

    // Create index file for plugins
    const indexContent = `---
title: Plugins
description: List of all available plugins and their documentation
---

# Available Plugins

${pluginDirs
  .map((dir) => {
    const name = normalizePath(dir).replace(/\/$/, '');
    return `- [${name}](${name})`;
  })
  .join('\n')}
`;

    await writeFile(join(pluginsDocsDir, 'index.md'), indexContent);
    console.log('Generated plugins index');
  } catch (error) {
    console.error('Error collecting plugins documentation:', error);
  }
}

async function collectDocs() {
  try {
    // Find all README.md files in the project
    const files = await glob('**/*.md', {
      cwd: rootDir,
      ignore: [...ignorePatterns, '**/connectors/**', '**/plugins/**'],
      nocase: true, // Case-insensitive search
    });

    console.log('Found files:', files);

    for (const file of files) {
      await copyFile(file);
    }

    // Copy assets
    await copyAssets();

    // Collect plugins and connectors documentation
    await collectPluginsDocs();
    await collectConnectorsDocs();

    console.log('Documentation collection completed!');
    return files;
  } catch (error) {
    console.error('Error collecting documentation:', error);
    process.exit(1);
  }
}

function shouldProcessFile(filepath) {
  // Check if file is README.md and not in ignored directories
  const isReadme = /readme\.md$/i.test(filepath);
  const isIgnored = ignorePatterns.some((pattern) => {
    const regexPattern = pattern.replace(/\*\*/g, '.*');
    return new RegExp(regexPattern, 'i').test(filepath);
  });

  return isReadme && !isIgnored;
}

async function cleanTargetDirs() {
  console.log('Cleaning target directories...');
  try {
    await rm(docsDir, { recursive: true, force: true });
    await rm(assetsDir, { recursive: true, force: true });
    console.log('Target directories cleaned successfully');
  } catch (error) {
    console.error('Error cleaning target directories:', error);
  }
}

async function watchFiles() {
  // Clean target directories first
  await cleanTargetDirs();

  // First collect all files
  const initialFiles = await collectDocs();

  console.log('Watching for file changes...', initialFiles);

  // Start watching for changes
  watch(rootDir, { recursive: true }, async (eventType, filename) => {
    if (!filename) return;

    // Normalize path for all platforms
    const relativePath = normalizePath(filename);

    if (shouldProcessFile(relativePath)) {
      console.log(`Change detected in ${relativePath}`);
      await copyFile(relativePath);
    }
  });
}

// Check command line arguments
const args = process.argv.slice(2);
if (args.includes('--watch')) {
  watchFiles();
} else {
  cleanTargetDirs().then(() => collectDocs());
}
