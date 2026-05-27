package x

import (
	"bufio"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	DirPerm  = 0755
	FilePerm = 0644
)

type Path struct {
	value string
}

func NewPath(elem ...string) Path {
	return Path{value: filepath.Join(elem...)}
}

func NewAbsPath(elem ...string) Path {
	abs, err := filepath.Abs(filepath.Join(elem...))
	if err != nil {
		return Path{}
	}
	return Path{value: abs}
}

// 当前程序路径
func Executable() Path {
	f, err := os.Executable()
	if err != nil {
		return Path{}
	}
	return NewPath(f)
}

func WorkingDir() Path {
	dir, err := os.Getwd()
	if err != nil {
		return Path{}
	}
	return NewPath(dir)
}

func HomeDir() Path {
	dir, err := os.UserHomeDir()
	if err != nil {
		return Path{}
	}
	return NewPath(dir)
}

func TempDir() Path {
	return NewPath(os.TempDir())
}

func (p Path) Len() int {
	return len(p.value)
}

func (p Path) String() string {
	return p.value
}

// 用 / 作为路径分隔符
func (p Path) Posix() string {
	return filepath.ToSlash(p.value)
}

// file:// 格式
func (p Path) Uri() string {
	abs := p.value
	if !filepath.IsAbs(abs) {
		if a, err := filepath.Abs(abs); err == nil {
			abs = a
		}
	}
	return "file://" + filepath.ToSlash(abs)
}

func (p Path) IsAbs() bool {
	return filepath.IsAbs(p.value)
}

func (p Path) Abs() Path {
	abs, err := filepath.Abs(p.value)
	if err != nil {
		return Path{}
	}
	return NewPath(abs)
}

// 文件名（带扩展名）
func (p Path) Name() string {
	return filepath.Base(p.value)
}

// 文件名（不带扩展名），
// abc -> abc，
// abc.pdf -> abc，
// abc.xy.docx -> abc.xy，
// .abc -> .abc，
// .abc.xy -> .abc
func (p Path) Stem() string {
	name := filepath.Base(p.value)
	i := strings.LastIndex(name, ".")
	if i <= 0 {
		return name
	}
	return name[:i]
}

// 扩展名（包括开头的点），
// abc -> 空，
// abc.pdf -> .pdf，
// abc.xy.docx -> .docx，
// .abc -> 空，
// .abc.xy -> .xy
func (p Path) Ext() string {
	name := filepath.Base(p.value)
	i := strings.LastIndex(name, ".")
	if i <= 0 {
		return ""
	}
	return name[i:]
}

func (p Path) WithName(name string) Path {
	return NewPath(filepath.Dir(p.value), name)
}

func (p Path) WithStem(stem string) Path {
	return NewPath(filepath.Dir(p.value), stem+p.Ext())
}

func (p Path) WithExt(ext string) Path {
	return NewPath(filepath.Dir(p.value), p.Stem()+ext)
}

// 盘符（仅 Windows，包括冒号），
// 对于 UNC 路径，是服务器名和共享名的组合（如 \\server\share）
func (p Path) Drive() string {
	return filepath.VolumeName(p.value)
}

// 如果是绝对路径，就返回路径分隔符，否则为空，
// UNC 路径 \\server\share 是绝对路径
func (p Path) Root() string {
	if filepath.IsAbs(p.value) {
		return PathSeparator
	}
	return ""
}

// Drive + Root
func (p Path) Anchor() string {
	return p.Drive() + p.Root()
}

func (p Path) Parent() Path {
	return NewPath(filepath.Dir(p.value))
}

// 相对路径的最后一个元素是点（当前目录），
// 例如 abc/xyz/abc.pdf -> [abc/xyz, abc, .]，
// 绝对路径的最后一个元素是根目录（Windows 是盘符包括冒号，其他是 /），
// 例如 /usr/local/abc -> [/usr/local, /usr, /]，
// C:\Windows\system32\abc.dll -> [C:\Windows\system32, C:\Windows, C:]，
// 在 Windows 上，以 / 开头的视为相对路径，
// 非 Windows 上，以盘符开头的视为相对路径
func (p Path) Parents() []Path {
	var parents []Path
	cur := p.value
	for {
		dir := filepath.Dir(cur)
		if dir == cur {
			break
		}
		parents = append(parents, NewPath(dir))
		cur = dir
	}
	return parents
}

// 把路径按分隔符分隔为单独的部分，没有表示当前目录的点，
// 非 Windows 上也没有根目录 /（但 Windows 上有盘符包括冒号）
func (p Path) Parts() []string {
	var parts []string
	for s := range strings.SplitSeq(p.value, PathSeparator) {
		if s != "" {
			parts = append(parts, s)
		}
	}
	return parts
}

func (p Path) Join(elem ...string) Path {
	return NewPath(append([]string{p.value}, elem...)...)
}

// 类似于 cd，但只是构造 Path，不是切换程序工作目录，
// 如果路径是文件，则返回文件所在目录 + <path>，
// 例如 abc/xyz/abc.pdf -> abc/xyz/<path>，
// 如果路径是目录，则返回自身 + <path>，
// 例如 abc/xyz/docs -> abc/xyz/docs/<path>，
// 如果出错（路径不存在或无权限等）则返回空
func (p Path) JoinDir(path string) Path {
	info := p.Stat()
	if info == nil {
		return Path{}
	}
	if info.IsDir() {
		return NewPath(p.value, path)
	}
	return NewPath(filepath.Dir(p.value), path)
}

func (p Path) Relative(base string) Path {
	rel, err := filepath.Rel(base, p.value)
	if err != nil {
		return Path{}
	}
	return NewPath(rel)
}

func (p Path) Matches(pattern string) bool {
	m, err := filepath.Match(pattern, p.value)
	return err == nil && m
}

func (p Path) Stat() fs.FileInfo {
	info, err := os.Stat(p.value)
	if err != nil {
		return nil
	}
	return info
}

// 路径不存在或无权限等均视为 false
func (p Path) Exists() bool {
	return p.Stat() != nil
}

// 路径不存在或无权限等均视为 false
func (p Path) IsDir() bool {
	info := p.Stat()
	if info == nil {
		return false
	}
	return info.IsDir()
}

// 路径不存在或无权限等均视为 false
func (p Path) IsFile() bool {
	info := p.Stat()
	if info == nil {
		return false
	}
	return info.Mode().IsRegular()
}

// 路径不存在或无权限等均视为 false
func (p Path) IsSymlink() bool {
	info := p.Stat()
	if info == nil {
		return false
	}
	return info.Mode()&fs.ModeSymlink != 0
}

// 找到符号链接的目标路径（非符号链接原样返回），
// 如果出错（路径不存在或无权限等）则返回空
func (p Path) EvalSymlinks() Path {
	real, err := filepath.EvalSymlinks(p.value)
	if err != nil {
		return Path{}
	}
	return NewPath(real)
}

func (p Path) Perm() fs.FileMode {
	info := p.Stat()
	if info == nil {
		return 0
	}
	return info.Mode().Perm()
}

func (p Path) Size() int64 {
	info := p.Stat()
	if info == nil {
		return -1
	}
	return info.Size()
}

// 判断目录或文件是否为空，
// 路径不存在或无权限等均视为 true
func (p Path) IsEmpty() bool {
	info := p.Stat()
	if info == nil {
		return true
	}
	if info.IsDir() {
		f, err := os.Open(p.value)
		if err != nil {
			return true
		}
		defer f.Close()
		_, err = f.Readdirnames(1)
		return err != nil
	}
	return info.Size() == 0
}

// 判断路径是否相同（展开符号链接），不比较内容
// 如果出错（路径不存在或无权限等）则返回 false
func (p Path) Equals(other string) bool {
	p1, err := filepath.Abs(p.value)
	if err != nil {
		return false
	}
	p2, err := filepath.Abs(other)
	if err != nil {
		return false
	}
	p1, err = filepath.EvalSymlinks(p1)
	if err != nil {
		return false
	}
	p2, err = filepath.EvalSymlinks(p2)
	if err != nil {
		return false
	}
	return p1 == p2
}

// 如果路径所指向的目录或文件存在，则打开，
// 如果不存在，则创建文件（自动创建不存在的目录）
// 如果出错（无权限等）则返回 nil
func (p Path) Open() *os.File {
	if err := os.MkdirAll(filepath.Dir(p.value), DirPerm); err == nil {
		if f, err := os.OpenFile(p.value, os.O_RDWR|os.O_CREATE, FilePerm); err == nil {
			return f
		}
	}
	return nil
}

// 如果路径存在且是目录，则返回 nil，
// 如果存在且是文件，则先清空，
// 如果不存在，则创建文件（自动创建不存在的目录），
// 如果出错（无权限等）则返回 nil
func (p Path) OpenWithTrunc() *os.File {
	if err := os.MkdirAll(filepath.Dir(p.value), DirPerm); err == nil {
		if f, err := os.OpenFile(p.value, os.O_RDWR|os.O_CREATE|os.O_TRUNC, FilePerm); err == nil {
			return f
		}
	}
	return nil
}

// 把 src 下的所有目录和文件都复制到 dest 下（不是复制 src 目录本身），
// src 和 dest 都是目录，已存在的文件会被覆盖
func (p Path) copyDir(dest, src string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil || path == src {
			return nil
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return nil
		}
		dst := filepath.Join(dest, rel)
		if !d.IsDir() {
			return p.copyFile(dst, path)
		}
		return nil
	})
}

// 把 src 复制到 dest，src 和 dest 都是文件，已存在则覆盖
func (p Path) copyFile(dest, src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	if err = os.MkdirAll(filepath.Dir(dest), DirPerm); err != nil {
		return err
	}
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
	return err
}

func (p Path) copy(dest string) error {
	srcInfo, err := os.Stat(p.value)
	if err != nil {
		return err
	}
	destInfo, err := os.Stat(dest)
	if srcInfo.IsDir() {
		if (err != nil && os.IsNotExist(err)) || (err == nil && destInfo.IsDir()) {
			return p.copyDir(dest, p.value)
		}
		return nil
	}
	if err != nil {
		if os.IsNotExist(err) {
			return p.copyFile(dest, p.value)
		}
		return nil
	}
	if destInfo.IsDir() {
		return p.copyFile(filepath.Join(dest, filepath.Base(p.value)), p.value)
	} else {
		return p.copyFile(dest, p.value)
	}
}

// 将 p 复制到 dest，
// 假设 p: /src/abc，dest: /dest/xyz，
// 如果 p 是目录，则将其复制到 dest 下，复制后目录结构为 /dest/xyz/abc，
// 如果 p 是文件，分三种情况：
//
//	1、如果 dest 存在且是目录，则将 p 复制到此目录下（文件名保持不变），复制后目录结构为 /dest/xyz/abc，
//	2、如果 dest 存在且是文件，则覆盖此文件，复制后目录结构为 /dest/xyz，xyz 就是 abc 文件，
//	3、如果 dest 不存在，则将 dest 视为文件路径，复制后目录结构为 /dest/xyz，xyz 就是 abc 文件
func (p Path) Copy(dest string) {
	p.copy(filepath.Join(dest, filepath.Base(p.value)))
}

// 跟 Copy 唯一的区别是复制的目录结构不一样，
// CopyAll 是把 p 下的所有子目录和文件复制到 dest 下，而 Copy 是直接把 p 本身复制到 dest 下
func (p Path) CopyAll(dest string) {
	p.copy(dest)
}

// 移动文件或目录，不支持跨文件系统（可以用复制+删除）
func (p Path) Move(dest string) {
	if err := os.MkdirAll(filepath.Dir(dest), DirPerm); err == nil {
		os.Rename(p.value, dest)
	}
}

// 重命名，要移动文件或目录使用 Move
func (p Path) Rename(newName string) {
	os.Rename(p.value, filepath.Join(filepath.Dir(p.value), newName))
}

func (p Path) Delete() {
	os.RemoveAll(p.value)
}

func (p Path) Symlink(target string) {
	if err := os.MkdirAll(filepath.Dir(p.value), DirPerm); err == nil {
		os.Symlink(target, p.value)
	}
}

func (p Path) Hardlink(target string) {
	if err := os.MkdirAll(filepath.Dir(p.value), DirPerm); err == nil {
		os.Link(target, p.value)
	}
}

func (p Path) Mkdirs() {
	os.MkdirAll(p.value, DirPerm)
}

func (p Path) Chmod(perm fs.FileMode) {
	os.Chmod(p.value, perm)
}

func (p Path) Chown(uid, gid int) {
	os.Chown(p.value, uid, gid)
}

// 修改访问时间和修改时间
func (p Path) Chtimes(atime, mtime time.Time) {
	os.Chtimes(p.value, atime, mtime)
}

// 一次性读取文件
func (p Path) ReadBytes() []byte {
	data, err := os.ReadFile(p.value)
	if err != nil {
		return nil
	}
	return data
}

// 逐块读取文件，bufLen 为缓冲区大小
func (p Path) ReadBytesBuf(bufLen int, fn func(buf []byte) error) {
	f, err := os.Open(p.value)
	if err != nil {
		return
	}
	defer f.Close()
	if bufLen <= 0 {
		bufLen = 32 * 1024
	}
	buf := make([]byte, bufLen)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF && n > 0 {
				fn(buf[:n])
			}
			break
		}
		if err = fn(buf[:n]); err != nil {
			break
		}
	}
}

// 一次性读取文件
func (p Path) ReadText() string {
	data, err := os.ReadFile(p.value)
	if err != nil {
		return ""
	}
	return string(data)
}

// 逐行读取文件，返回所有行
func (p Path) ReadLines() []string {
	f, err := os.Open(p.value)
	if err != nil {
		return nil
	}
	defer f.Close()
	lines := make([]string, 0)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if err = sc.Err(); err != nil {
		return nil
	}
	return lines
}

// 逐行读取文件
func (p Path) ReadLinesFunc(fn func(line string) error) {
	f, err := os.Open(p.value)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if err = fn(sc.Text()); err != nil {
			break
		}
	}
	_ = sc.Err()
}

// 读取 JSON，文件不存在或无权限等返回 nil
func (p Path) ReadJson() any {
	data, err := os.ReadFile(p.value)
	if err != nil {
		return nil
	}
	var a any
	if err = json.Unmarshal(data, &a); err != nil {
		return nil
	}
	return a
}

// 写入字节，文件存在则覆盖，不存在则创建
func (p Path) WriteBytes(data []byte) {
	if err := os.MkdirAll(filepath.Dir(p.value), DirPerm); err == nil {
		os.WriteFile(p.value, data, FilePerm)
	}
}

// 写入文本，文件存在则覆盖，不存在则创建
func (p Path) WriteText(text string) {
	if err := os.MkdirAll(filepath.Dir(p.value), DirPerm); err == nil {
		os.WriteFile(p.value, []byte(text), FilePerm)
	}
}

// 写入 JSON，文件存在则覆盖，不存在则创建
func (p Path) WriteJson(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(p.value), DirPerm); err == nil {
		os.WriteFile(p.value, data, FilePerm)
	}
}

// 追加字节，文件存在则覆盖，不存在则创建
func (p Path) AppendBytes(data []byte) {
	if err := os.MkdirAll(filepath.Dir(p.value), DirPerm); err != nil {
		return
	}
	f, err := os.OpenFile(p.value, os.O_WRONLY|os.O_CREATE|os.O_APPEND, FilePerm)
	if err != nil {
		return
	}
	defer f.Close()
	f.Write(data)
}

// 追加文本，文件存在则覆盖，不存在则创建
func (p Path) AppendText(text string) {
	p.AppendBytes([]byte(text))
}

// 返回直接子目录和文件
func (p Path) List() []string {
	entries, err := os.ReadDir(p.value)
	if err != nil {
		return nil
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

// 返回直接子目录和文件
func (p Path) ListFunc(fn func(string, fs.DirEntry) bool) []string {
	entries, err := os.ReadDir(p.value)
	if err != nil {
		return nil
	}
	var names []string
	for _, entry := range entries {
		if fn(filepath.Join(p.value, entry.Name()), entry) {
			names = append(names, entry.Name())
		}
	}
	return names
}

// 返回所有子目录和文件
func (p Path) ListAll() []string {
	var paths []string
	err := filepath.WalkDir(p.value, func(path string, _ fs.DirEntry, err error) error {
		if err != nil || path == p.value {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return nil
	}
	return paths
}

// 返回所有子目录和文件
func (p Path) ListAllFunc(fn func(string, fs.DirEntry) bool) []string {
	var paths []string
	err := filepath.WalkDir(p.value, func(path string, d fs.DirEntry, err error) error {
		if err != nil || path == p.value {
			return nil
		}
		if fn(path, d) {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return paths
}

// 递归遍历目录
func (p Path) Walk(fn fs.WalkDirFunc) {
	filepath.WalkDir(p.value, fn)
}
