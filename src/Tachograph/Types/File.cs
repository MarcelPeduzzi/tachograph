using System;

namespace Tachograph
{
    /// <summary>
    /// Represents the type of parsed tachograph file.
    /// </summary>
    public enum FileType
    {
        Unknown = 0,
        DriverCard = 1,
        VehicleUnit = 2,
        WorkshopCard = 3,
        ControlCard = 4,
        CompanyCard = 5
    }

    /// <summary>
    /// Represents a parsed tachograph file with semantic data structures.
    /// </summary>
    public class File
    {
        /// <summary>
        /// Gets or sets the type of the file.
        /// </summary>
        public FileType Type { get; set; }

        /// <summary>
        /// Gets or sets the driver card file data (when Type is DriverCard).
        /// </summary>
        public DriverCardFile DriverCard { get; set; }

        /// <summary>
        /// Gets or sets the vehicle unit file data (when Type is VehicleUnit).
        /// </summary>
        public VehicleUnitFile VehicleUnit { get; set; }

        /// <summary>
        /// Gets or sets the workshop card file data (when Type is WorkshopCard).
        /// </summary>
        public WorkshopCardFile WorkshopCard { get; set; }

        /// <summary>
        /// Gets or sets the control card file data (when Type is ControlCard).
        /// </summary>
        public ControlCardFile ControlCard { get; set; }

        /// <summary>
        /// Gets or sets the company card file data (when Type is CompanyCard).
        /// </summary>
        public CompanyCardFile CompanyCard { get; set; }
    }
}
