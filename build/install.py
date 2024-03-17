#!/usr/bin/env python3

import os
import shutil
import subprocess
import sys
from pathlib import Path


class FileManager:
    @staticmethod
    def set_permissions_and_make_removable(target_dir="target", file_permission=0o755, dir_permission=0o777):
        target_dir_path = Path(target_dir)
        for path in target_dir_path.glob('**/*'):
            os.chmod(path, file_permission if path.is_file() else dir_permission)

    @staticmethod
    def move_executable(build_executable_path, target_dir, executable_name):
        target_path = Path(target_dir) / 'theHamdiz' / executable_name
        if not target_path.parent.exists():
            target_path.parent.mkdir(parents=True)
        shutil.move(str(build_executable_path), str(target_path))
        return target_path

    @staticmethod
    def create_symlink(source_path, symlink_path):
        if symlink_path.exists():
            if symlink_path.is_symlink() or symlink_path.is_file():
                symlink_path.unlink()  # Removes the symlink or file if it exists
            elif symlink_path.is_dir():
                shutil.rmtree(symlink_path)  # Removes the directory if it's a symlink pointing to a directory

        source_path = Path(source_path)  # Ensure source_path is a Path object for consistency
        symlink_path.symlink_to(source_path, target_is_directory=source_path.is_dir())


class EnvironmentManager:
    @staticmethod
    def add_to_system_path(target_path):
        try:
            # Get the absolute path of the target directory
            target_dir_path = Path(target_path).resolve()
            # Get the current PATH environment variable
            current_path = os.environ.get('PATH', '')
            # Check if the target directory is already in PATH
            if str(target_dir_path) not in current_path:
                # If it's not, append it to PATH
                new_path = f"{current_path};{target_dir_path}"
                # Set the new PATH environment variable
                subprocess.run(['setx', 'PATH', new_path], check=True)
                print('> Successfully added the target directory to PATH.')
        except Exception as e:
            print(
                f'> Error: Could not add the target directory to PATH.\n> {e}\n> Try running the script as an '
                f'administrator.')
            sys.exit(1)


class BuildManager:
    def __init__(self, platform=sys.platform):
        self.platform = platform
        self.build_scripts = {
            'win32': 'build\\build.bat',
            'linux': 'build/build.sh',
            'darwin': 'build/build.sh'
        }
        self.target_dirs = {
            'win32': 'C:\\Program Files\\',
            'linux': '/usr/local/bin/',
            'darwin': '/usr/local/bin/'
        }
        self.executable = 'nsh.exe' if platform == 'win32' else 'nsh'

    def run_build_script(self):
        build_script = self.build_scripts.get(self.platform)
        if build_script:
            try:
                subprocess.run(['./' + build_script], shell=True, check=True)
                print(f'> Successfully built the project using {build_script}.')
            except subprocess.CalledProcessError as e:
                print(f'> Error: Could not find or execute the build script {build_script}.\n> {e}')
                sys.exit(1)

    def install_executable(self):
        FileManager.set_permissions_and_make_removable()
        print(f'> Successfully made {self.executable} executable.')

        target_dir = self.target_dirs.get(self.platform)
        if target_dir:
            build_executable_path = Path('target', self.platform, self.executable).resolve()
            target_path = FileManager.move_executable(build_executable_path, target_dir, self.executable)
            print(f'> Successfully installed {self.executable} to {target_path}.')

            if self.platform != 'win32':
                symlink_path = Path('/usr/local/bin', self.executable)
                FileManager.create_symlink(target_path, symlink_path)
                print(
                    '> Successfully created a symlink to the executable system-wide, you can now run `nsh` from '
                    'anywhere.')
                os.chmod(target_path, 0o755)
                print('> Successfully set the permissions of the executable files.')
            else:
                EnvironmentManager.add_to_system_path(target_path)


if __name__ == '__main__':
    build_manager = BuildManager()
    build_manager.run_build_script()
    build_manager.install_executable()
